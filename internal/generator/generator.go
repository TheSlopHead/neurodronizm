package generator

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/genai"
)

type Generator struct {
	genaiClient *genai.Client
	apiKey      string
}

func New(ctx context.Context, apiKey string) (*Generator, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("Cannot create gemini client: %v", err)
	}
	return &Generator{
		genaiClient: client,
		apiKey:      apiKey,
	}, nil
}

func (g *Generator) GeneratePost(ctx context.Context, examples []string, topic string) ([]string, error) {
	Prompt := strings.Join(examples, "\n---\n")
	finalPrompt := fmt.Sprintf("Вот примеры моих постов:\n%s\n\nА теперь: %s", Prompt, topic)

	parts := []*genai.Part{
		{Text: finalPrompt},
	}
	rawResult, err := g.genaiClient.Models.GenerateContent(ctx, "gemini-3.1-flash-lite", []*genai.Content{{Parts: parts}}, nil)
	if err != nil {
		return nil, fmt.Errorf("Cannot generate post: %v", err)
	}
	res := rawResult.Text()
	variants := strings.Split(res, "[POST_SPLIT]")
	cleanVariants := make([]string, 0, len(variants))

	for _, variant := range variants {
		trimmed := strings.TrimSpace(variant)
		if trimmed != "" {
			cleanVariants = append(cleanVariants, trimmed)
		}
	}
	return cleanVariants, nil
}
