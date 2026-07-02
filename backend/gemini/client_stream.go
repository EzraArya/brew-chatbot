package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

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

	config := GetToolConfig()

	chat, err := c.ai.Chats.Create(ctx, "gemini-2.5-flash", config, geminiHistory)

	if err != nil {
		return fmt.Errorf("failed to create chat session: %w", err)
	}

	stream := chat.SendMessageStream(ctx, genai.Part{Text: userMessage})

	for resp, err := range stream {
		if err != nil {
			return fmt.Errorf("stream failed: %w", err)
		}

		if functionCalls := resp.FunctionCalls(); len(functionCalls) > 0 {
			call := functionCalls[0]

			jsonBytes, err := json.Marshal(call.Args)
			if err == nil {
				slog.Info("gemini executed tool", "tool", call.Name, "args", string(jsonBytes))
				onChunk(fmt.Sprintf("[%s] %s", call.Name, string(jsonBytes)))
			}

			return nil
		}
		
        if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
        	chunkText := resp.Candidates[0].Content.Parts[0].Text
         	if chunkText != "" {
          		slog.Debug("gemini generated text", "length", len(chunkText))
         		onChunk(chunkText)
         	}
        }
	}

	return nil
}