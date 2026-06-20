package gemini

import (
	"context"
	"fmt"

	"google.golang.org/genai"
)

// The system prompt — this is the "personality" of your chatbot
const systemPrompt = `You are BrewBot, a friendly and knowledgeable brewing expert. 
You help users with all aspects of brewing including:
- Beer, coffee, tea, and kombucha brewing
- Recipes, ingredients, and techniques
- Troubleshooting brewing problems
- Equipment recommendations

Keep your answers practical, friendly, and concise.
If a question is not related to brewing, politely redirect the conversation back to brewing topics.`

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
	chat, err := c.ai.Chats.Create(ctx, "gemini-2.5-flash", &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{{Text: systemPrompt}},
		},
	}, geminiHistory)
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