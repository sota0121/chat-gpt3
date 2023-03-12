package application

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/sashabaranov/go-openai"
	"github.com/sota0121/go-ai-chat/internal"
	"golang.org/x/exp/slog"
)

type ChatService interface {
	SendText(ctx context.Context, text string) error
	SendTextStream(ctx context.Context, text string) error
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

func (s *chatService) SendTextStream(ctx context.Context, text string) error {
	openaiClient := internal.GetOpenAIClientFromContext(ctx)

	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: text,
			},
		},
		Stream: true,
	}

	stream, err := openaiClient.CreateChatCompletionStream(ctx, req)
	if err != nil {
		slog.Error("Error creating chat completion stream", err)
		return err
	}
	defer stream.Close()

	fmt.Printf("AI> ")
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Printf("\n\n")
			return nil
		}
		if err != nil {
			slog.Error("Error receiving chat completion stream", err)
			return err
		}

		if len(response.Choices) == 0 {
			continue
		}
		fmt.Printf("%v", response.Choices[0].Delta.Content)
	}
}
