package gemini

import (
	"context"
	"fmt"

	"google.golang.org/genai"
)

// The system prompt — this is the "personality" of your chatbot
const systemPrompt = `You are BrewBot, a friendly and knowledgeable brewing expert.
You help users with all aspects of brewing including beer, coffee, tea, and kombucha.

## Tool Usage Guidelines
Use tools when the response would benefit from structured, interactive display.
Use plain text for general questions, advice, and explanations.

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
// json tags tell Go how to convert this to/from JSON
type Message struct {
	Role    string `json:"role"`    // "user" or "model"
	Content string `json:"content"`
}

// Client wraps the Gemini SDK client
type Client struct {
	ai *genai.Client
}

// NewClient creates and returns a new Gemini client
// This is Go's equivalent of an initializer/constructor
func NewClient(apiKey string) (*Client, error) {
	ctx := context.Background()

	ai, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	return &Client{ai: ai}, nil
}

// Chat sends a message to Gemini with the full conversation history
// and returns the assistant's reply
func (c *Client) Chat(ctx context.Context, history []Message, userMessage string) (string, error) {
	// Convert our Message slice into Gemini's format
	var geminiHistory []*genai.Content
	for _, msg := range history {
		geminiHistory = append(geminiHistory, &genai.Content{
			Role:  msg.Role,
			Parts: []*genai.Part{{Text: msg.Content}},
		})
	}

	// Create a chat session with history + system prompt
	config := GetToolConfig()

	chat, err := c.ai.Chats.Create(ctx, "gemini-2.5-flash", config, geminiHistory)
	if err != nil {
		return "", fmt.Errorf("failed to create chat session: %w", err)
	}

	resp, err := chat.SendMessage(ctx, genai.Part{Text: userMessage})
	if err != nil {
		return "", fmt.Errorf("failed to send message: %w", err)
	}

	// Extract the text response
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("empty response from Gemini")
	}

	return resp.Candidates[0].Content.Parts[0].Text, nil
}