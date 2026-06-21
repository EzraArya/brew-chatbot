package gemini

import (
	"context"
	"fmt"

	"google.golang.org/genai"
)

func (c *Client) ChatStream(
	ctx context.Context,
	history []Message,
	userMessage string,
	onChunk func( chunk string),
) error {
	var geminiHistory []*genai.Content
	for _, msg := range history {
		geminiHistory = append(geminiHistory, &genai.Content{
			Role:    msg.Role,
			Parts:   []*genai.Part{{Text: msg.Content}},
		})
	}

	chat, err := c.ai.Chats.Create(ctx, "gemini-2.5-flash", &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{{Text: systemPrompt}},
		},
	}, geminiHistory)

	if err != nil {
		return fmt.Errorf("failed to create chat session: %w", err)
	}

	stream := chat.SendMessageStream(ctx, genai.Part{Text: userMessage})

	for resp, err := range stream {
		if err != nil {
			return fmt.Errorf("stream failed: %w", err)
		}
        if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
        	chunkText := resp.Candidates[0].Content.Parts[0].Text
         	if chunkText != "" {
         		onChunk(chunkText)
         	}
        }
	}

	return nil
}