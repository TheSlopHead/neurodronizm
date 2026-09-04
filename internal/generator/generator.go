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

// type Generator struct {
// 	openaiClient *openai.Client
// 	apiKey       string
// }

// func New(ctx context.Context, apiKey string) (*Generator, error) {
// 	client := openai.NewClient(
// 		option.WithAPIKey("ollama"),
// 		option.WithBaseURL("http://localhost:11434/v1/"),
// 	)
// 	return &Generator{
// 		openaiClient: &client,
// 		apiKey:       " ",
// 	}, nil

// }

// func (g *Generator) GeneratePost(ctx context.Context, examples []string, topic string) ([]string, error) {
// 	Prompt := strings.Join(examples, "\n---\n")
// 	finalPrompt := fmt.Sprintf("Вот примеры моих постов:\n%s\n\nА теперь: %s", Prompt, topic)

// 	rawResult, err := g.openaiClient.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
// 		Model: openai.ChatModel("llama3.2:3b"),
// 		Messages: ([]openai.ChatCompletionMessageParamUnion{
// 			openai.UserMessage(finalPrompt),
// 		}),
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("Cannot generate post: %v", err)
// 	}

// 	if len(rawResult.Choices) == 0 {
// 		return nil, fmt.Errorf("Ollama returned empty choices")
// 	}
// 	res := rawResult.Choices[0].Message.Content
// 	variants := strings.Split(res, "[POST_SPLIT]")
// 	cleanVariants := make([]string, 0, len(variants))

// 	for _, variant := range variants {
// 		trimmed := strings.TrimSpace(variant)
// 		if trimmed != "" {
// 			cleanVariants = append(cleanVariants, trimmed)
// 		}
// 	}
// 	return cleanVariants, nil
// }

// func (g *Generator) GetEmbedding(ctx context.Context, text string) ([]float32, error) {
// 	rawEmbed, err := g.openaiClient.Embeddings.New(ctx,openai.EmbeddingNewParams{
// 			Model: openai.ChatModel("nomic-embed-text"),
// 			Input: (openai.EmbeddingNewParamsInputUnion{
// 				OfString: openai.String(text),
// 			}),
// 		})
// 	if err != nil {
// 		return nil, fmt.Errorf("ollama embedding failed: %v", err)
// 	}

// 	vec64 := rawEmbed.Data[0].Embedding
// 	vec32 := make([]float32, len(vec64))
// 	for i, v := range vec64 {
// 		vec32[i] = float32(v)
// 	}
// 	return vec32, nil
// }
