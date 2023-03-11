package internal

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

type OpenAIClientKey struct{}

// SetOpenAIClientToContext sets OpenAI client to context
func SetOpenAIClientToContext(ctx context.Context, client *openai.Client) context.Context {
	return context.WithValue(ctx, OpenAIClientKey{}, client)
}

// GetOpenAIClientFromContext gets OpenAI client from context
func GetOpenAIClientFromContext(ctx context.Context) *openai.Client {
	return ctx.Value(OpenAIClientKey{}).(*openai.Client)
}
