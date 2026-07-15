package gemini

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/genai"
)

// GenerateTitle produces a short, specific conversation title based on
// the user's first message. It is intentionally a one-shot call with no
// chat history or tools — just a focused prompt for a short title.
func (c *Client) GenerateTitle(ctx context.Context, userMessage string) (string, error) {
	prompt := fmt.Sprintf(`Based on this brewing question:
"%s"

Generate a short, specific title for this conversation (maximum 5 words).
Rules: No punctuation. No quotes. Just the title itself.`, userMessage)

	temp := float32(0.3)

	resp, err := c.ai.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(prompt),
		&genai.GenerateContentConfig{
			SystemInstruction: &genai.Content{
				Parts: []*genai.Part{{Text: "You generate short, specific conversation titles. Reply with only the title — nothing else."}},
			},
			Temperature:      &temp,
		},
	)
	if err != nil {
		return "", fmt.Errorf("generating title: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("empty title response")
	}

	title := strings.TrimSpace(resp.Candidates[0].Content.Parts[0].Text)
	return title, nil
}