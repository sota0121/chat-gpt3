package application

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/sashabaranov/go-openai"
	"github.com/sota0121/go-ai-chat/internal"
	"golang.org/x/exp/slog"
)

type ChatService interface {
	SendText(ctx context.Context, text string) error
	SendTextStream(ctx context.Context, text string) error
}

func NewChatService(systemMessages []string) ChatService {
	systemMessageContents := make([]openai.ChatCompletionMessage, 0, len(systemMessages))

	for _, msg := range systemMessages {
		systemMessageContents = append(systemMessageContents, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: msg,
		})
	}

	return &chatService{
		SystemMessages: systemMessageContents,
	}
}

type chatService struct {
	SystemMessages []openai.ChatCompletionMessage
	Histories      []openai.ChatCompletionMessage
}

var _ ChatService = (*chatService)(nil)

func (s *chatService) SendText(ctx context.Context, text string) error {
	openaiClient := internal.GetOpenAIClientFromContext(ctx)

	s.Histories = append(s.Histories, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: text,
	})

	allMessages := append(s.SystemMessages, s.Histories...)

	req := openai.ChatCompletionRequest{
		Model:    openai.GPT3Dot5Turbo,
		Messages: allMessages,
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

	asistantResponse := response.Choices[0].Message
	s.Histories = append(s.Histories, asistantResponse)

	fmt.Printf("AI> %v\n", asistantResponse.Content)
	return nil
}

func (s *chatService) SendTextStream(ctx context.Context, text string) error {
	openaiClient := internal.GetOpenAIClientFromContext(ctx)

	s.Histories = append(s.Histories, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: text,
	})

	allMessages := append(s.SystemMessages, s.Histories...)

	req := openai.ChatCompletionRequest{
		Model:    openai.GPT3Dot5Turbo,
		Messages: allMessages,
		Stream:   true,
	}

	stream, err := openaiClient.CreateChatCompletionStream(ctx, req)
	if err != nil {
		slog.Error("Error creating chat completion stream", err)
		return err
	}
	defer stream.Close()

	responseTokens := []string{}
	fmt.Printf("AI> ")
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Printf("\n\n")
			asistantResponse := openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: strings.Join(responseTokens, ""),
			}
			s.Histories = append(s.Histories, asistantResponse)
			return nil
		}
		if err != nil {
			slog.Error("Error receiving chat completion stream", err)
			return err
		}

		if len(response.Choices) == 0 {
			slog.Error("Error creating chat completion", errors.New("no choices"))
			continue
		}

		responseTokens = append(responseTokens, response.Choices[0].Delta.Content)
		fmt.Printf("%v", response.Choices[0].Delta.Content)
	}
}
