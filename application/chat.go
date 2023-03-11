package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/sashabaranov/go-openai"
	"github.com/sota0121/go-ai-chat/internal"
	"golang.org/x/exp/slog"
)

type ChatService interface {
	SendText(ctx context.Context, text string) error
}

func NewChatService() ChatService {
	return &chatService{}
}

type chatService struct {
}

var _ ChatService = (*chatService)(nil)

func (s *chatService) SendText(ctx context.Context, text string) error {
	openaiClient := internal.GetOpenAIClientFromContext(ctx)

	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: text,
			},
		},
	}

	response, err := openaiClient.CreateChatCompletion(ctx, req)
	if err != nil {
		slog.Error("Error creating chat completion", err)
		return err
	}

	if len(response.Choices) == 0 {
		slog.Error("Error creating chat completion", errors.New("no choices"))
		return errors.New("no choices")
	}

	fmt.Printf("AI> %v\n", response.Choices[0].Message.Content)
	return nil
}
