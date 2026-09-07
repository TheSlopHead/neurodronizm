package generator

import (
	"context"
	"fmt"
	"strings"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type Generator struct {
	openaiClient *openai.Client
	apiKey       string
}

func New(ctx context.Context, apiKey string) (*Generator, error) {
	client := openai.NewClient(
		option.WithAPIKey("ollama"),
		option.WithBaseURL("http://localhost:11434/v1/"),
	)
	return &Generator{
		openaiClient: &client,
		apiKey:       " ",
	}, nil

}

func (g *Generator) GeneratePost(ctx context.Context, examples []string, topic string) ([]string, error) {
	Prompt := strings.Join(examples, "\n---\n")
	finalPrompt := fmt.Sprintf(`Вот примеры моих прошлых постов для понимания стиля:
	<переменная со склеенными примерами>

	Твоя задача: напиши 3 разных варианта нового поста на тему: "<переменная темы>".
	Пиши в моем стиле, иронично и со стебом.

	КРИТИЧЕСКИ ВАЖНОЕ ПРАВИЛО ДЛЯ ФОРМАТИРОВАНИЯ:
	Разделяй варианты постов строго строкой [POST_SPLIT].

	Внутри самих текстов постов этот маркер использовать запрещено. Не пиши никаких вступлений от себя вроде "Вот твои посты:".
	`, Prompt, topic)

	rawResult, err := g.openaiClient.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: openai.ChatModel("qwen2.5:1.5b"),
		Messages: ([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(finalPrompt),
		}),
	})
	if err != nil {
		return nil, fmt.Errorf("Cannot generate post: %v", err)
	}

	if len(rawResult.Choices) == 0 {
		return nil, fmt.Errorf("Ollama returned empty choices")
	}
	res := rawResult.Choices[0].Message.Content
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

func (g *Generator) GetEmbedding(ctx context.Context, text string) ([]float32, error) {
	rawEmbed, err := g.openaiClient.Embeddings.New(ctx, openai.EmbeddingNewParams{
		Model: openai.ChatModel("nomic-embed-text"),
		Input: (openai.EmbeddingNewParamsInputUnion{
			OfString: openai.String(text),
		}),
	})
	if err != nil {
		return nil, fmt.Errorf("ollama embedding failed: %v", err)
	}

	vec64 := rawEmbed.Data[0].Embedding
	vec32 := make([]float32, len(vec64))
	for i, v := range vec64 {
		vec32[i] = float32(v)
	}
	return vec32, nil
}

func (g *Generator) TopicGenerator(ctx context.Context, topic string) (string, error) {
	if topic == "" {
		prompt := "Ты — креативный автор канала. Придумай одну интересную, актуальную и острую тему для постироничного поста . Верни ТОЛЬКО саму тему одной короткой фразой, без лишних слов, кавычек и приветствий"
		rawResult, err := g.openaiClient.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
			Model: openai.ChatModel("qwen2.5:1.5b"),
			Messages: ([]openai.ChatCompletionMessageParamUnion{
				openai.UserMessage(prompt),
			}),
		})
		if err != nil {
			return "", fmt.Errorf("Cannot generate post: %v", err)
		}
		res := rawResult.Choices[0].Message.Content
		return res, nil
	} else {
		return topic, nil
	}

}
