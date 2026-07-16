package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"brew-chatbot/internal/rag"
	"google.golang.org/genai"
)

// The system prompt — this is the "personality" of your chatbot
const systemPrompt = `You are BrewBot, a friendly and knowledgeable brewing expert.
You help users with all aspects of brewing including beer, coffee, tea, and kombucha.

## Tool Usage Guidelines
Use tools when the response would benefit from structured, interactive display.
Use plain text for general questions, advice, and explanations.

- search_knowledge_base: ALWAYS call this first before answering any technical brewing question.
  Use it to look up precise parameters, ratios, water chemistry, grind sizes, or troubleshooting steps.
- generate_brew_recipe: For coffee manual brew recipes (V60, Chemex, AeroPress, French Press, etc.)
- generate_beer_recipe: For homebrewing beer recipes (IPA, stout, wheat beer, etc.)
- generate_tea_recipe: For tea preparation with specific parameters
- generate_kombucha_recipe: For kombucha fermentation recipes
- generate_troubleshooting: When the user describes a problem with their brew (bitter, sour, no carbonation, etc.)
- generate_brew_timer: ONLY when the user wants to brew RIGHT NOW and needs a real-time countdown.
  This is distinct from a recipe — timers have step-by-step durations for active brewing guidance.

## Behaviour
Keep answers practical, friendly, and concise.
If a question is unrelated to brewing, politely redirect back to brewing topics.
When a user asks for a recipe or timer, call the appropriate tool immediately using
sensible defaults. Do not ask clarifying questions before calling a tool — generate
the recipe first, then offer to adjust it afterwards.`

// Message represents a single chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Client wraps the Gemini SDK client and knowledge retriever
type Client struct {
	ai        *genai.Client
	retriever rag.Retriever
}

// NewClient creates a Gemini client wired with the markdown knowledge base
func NewClient(apiKey string) (*Client, error) {
	ctx := context.Background()

	ai, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	// Wire up RAG: markdown files → full context retriever
	source := &rag.MarkdownSource{Dir: "knowledge/"}
	retriever := &rag.FullContextRetriever{Sources: []rag.KnowledgeSource{source}}

	return &Client{ai: ai, retriever: retriever}, nil
}

// ChatStream sends a message to Gemini with conversation history and streams
// the response back chunk by chunk via the onChunk callback.
//
// search_knowledge_base tool calls are handled internally — the backend
// retrieves the docs and sends them back to Gemini transparently.
// All other tool calls (brew_recipe, timer, etc.) are forwarded to the caller via onChunk.
func (c *Client) ChatStream(
	ctx context.Context,
	history []Message,
	userMessage string,
	onChunk func(chunk string),
) error {
	var geminiHistory []*genai.Content
	for _, msg := range history {
		geminiHistory = append(geminiHistory, &genai.Content{
			Role:  msg.Role,
			Parts: []*genai.Part{{Text: msg.Content}},
		})
	}

	chat, err := c.ai.Chats.Create(ctx, "gemini-2.5-flash", GetToolConfig(), geminiHistory)
	if err != nil {
		return fmt.Errorf("failed to create chat session: %w", err)
	}

	stream := chat.SendMessageStream(ctx, genai.Part{Text: userMessage})

	// maxRAGRounds prevents infinite loops if Gemini keeps calling search_knowledge_base
	const maxRAGRounds = 3

	for round := 0; round < maxRAGRounds; round++ {
		var ragQuery string // set when Gemini calls search_knowledge_base

		for resp, err := range stream {
			if err != nil {
				return fmt.Errorf("stream failed: %w", err)
			}

			if functionCalls := resp.FunctionCalls(); len(functionCalls) > 0 {
				call := functionCalls[0]

				if call.Name == "search_knowledge_base" {
					// Record the query and use continue — do NOT break.
					// Breaking stops the Go iterator early, making a second
					// for-range on the same stream a no-op. The SDK only records
					// the function call in chat history after the stream is fully
					// consumed. Continuing lets it drain naturally.
					ragQuery, _ = call.Args["query"].(string)
					slog.Info("rag: knowledge search", "query", ragQuery)
					continue
				}

				// All other tools are forwarded to the iOS client
				jsonBytes, err := json.Marshal(call.Args)
				if err == nil {
					slog.Info("gemini executed tool", "tool", call.Name, "args", string(jsonBytes))
					onChunk(fmt.Sprintf("[%s] %s", call.Name, string(jsonBytes)))
				}

				return nil
			}

			// Only forward text chunks when not in RAG drain mode
			if ragQuery == "" {
				if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
					chunkText := resp.Candidates[0].Content.Parts[0].Text
					if chunkText != "" {
						slog.Debug("gemini generated text", "length", len(chunkText))
						onChunk(chunkText)
					}
				}
			}
		}

		// Stream fully consumed — SDK has recorded the model turn in history.
		// If no RAG call happened, we're done.
		if ragQuery == "" {
			break
		}

		// Retrieve knowledge and send function response back to Gemini
		docs, err := c.retriever.Retrieve(ctx, ragQuery)
		if err != nil {
			slog.Warn("rag: retrieval failed", "error", err)
			docs = nil
		}

		stream = chat.SendMessageStream(ctx, genai.Part{
			FunctionResponse: &genai.FunctionResponse{
				Name:     "search_knowledge_base",
				Response: map[string]any{"knowledge": buildKnowledgeResult(docs)},
			},
		})
	}

	return nil
}

// buildKnowledgeResult formats retrieved documents into a string for Gemini
func buildKnowledgeResult(docs []rag.Document) string {
	if len(docs) == 0 {
		return "No knowledge found for this query."
	}

	var sb strings.Builder
	for _, doc := range docs {
		sb.WriteString("## ")
		sb.WriteString(doc.Title)
		sb.WriteString("\n\n")
		sb.WriteString(doc.Content)
		sb.WriteString("\n\n---\n\n")
	}
	return sb.String()
}