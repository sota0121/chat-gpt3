package application

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
	"github.com/sota0121/go-ai-chat/internal"
	"golang.org/x/exp/slog"
)

type FindBugService interface {
	SendRequestStream(ctx context.Context, text string) error
}

func NewFindBugService() FindBugService {
	return &findBugService{}
}

type findBugService struct {
}

var _ FindBugService = (*findBugService)(nil)

const (
	findBugsMessageHeader = `以下のプログラムについて、バグを見つけてください。
	プログラムの関数名から、関数の満たすべき仕様を読み取ってください。
	バグを見つけたら、バグの原因と修正方法を説明してください。
	`
)

// SendRequestStream sends request to OpenAI to find bugs in stream
// This expects text to be in the following format:
// :findbugs <file> or :testgen <file> <function>
func (s *findBugService) SendRequestStream(ctx context.Context, text string) error {
	// Parse input text
	fileName, funcName, err := s.parseInput(text)
	if err != nil {
		return err
	}

	// Read file
	file, err := os.Open(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("file not found")
		}
		return err
	}

	// Get program code
	code, err := extractCode(file, funcName)
	if err != nil {
		slog.Error("Error extracting code", err)
		return err
	}

	// Make message body
	messageBody := fmt.Sprintf("%s\n\n%s", findBugsMessageHeader, code)

	// Get OpenAI client from context
	openaiClient := internal.GetOpenAIClientFromContext(ctx)

	// Send request to OpenAI in stream
	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: messageBody,
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

// parseInput parses input text
// [FIXME] This method is duplicated in test_gen.go
func (s *findBugService) parseInput(text string) (string, string, error) {
	// Check if text is in the correct format
	if err := s.checkInput(text); err != nil {
		return "", "", err
	}

	// Parse text
	tokens := strings.Split(text, " ")
	fileName := tokens[1]
	funcName := ""
	if len(tokens) > 2 {
		funcName = tokens[2]
	}
	return fileName, funcName, nil
}

// checkInput checks input text
func (s *findBugService) checkInput(text string) error {
	if !strings.HasPrefix(text, ":findbugs") {
		return errors.New("invalid format: text must start with ':findbugs'")
	}
	tokens := strings.Split(text, " ")
	if len(tokens) < 2 {
		return errors.New("invalid format: text must contain file name")
	}
	return nil
}
