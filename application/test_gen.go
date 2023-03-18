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

type TestGenService interface {
	SendRequest(ctx context.Context, text string) error
	SendRequestStream(ctx context.Context, text string) error
}

func NewTestGenService() TestGenService {
	return &testGenService{}
}

type testGenService struct {
}

var _ TestGenService = (*testGenService)(nil)

const (
	tesgGenMessageHeader = `以下のプログラムについて、テストコードを生成してください。
	テスト対象の関数の正常系と異常系のテストコードを生成してください。
	テストコード生成には、 gomock を使用してください。
	テストコードの形式は、AAA(Arrange, Act, Assert) に従ってください。
	`
)

// SendRequest sends request to OpenAI to generate test code
// This expects text to be in the following format:
// :testgen <file> or :testgen <file> <function>
func (s *testGenService) SendRequest(ctx context.Context, text string) error {
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
	messageBody := fmt.Sprintf("%s\n\n%s", tesgGenMessageHeader, code)

	// Get OpenAI client from context
	openaiClient := internal.GetOpenAIClientFromContext(ctx)

	// Send request to OpenAI
	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: messageBody,
			},
		},
	}
	response, err := openaiClient.CreateChatCompletion(ctx, req)
	if err != nil {
		slog.Error("Error creating chat completion", err)
		return err
	}

	// Print response
	if len(response.Choices) == 0 {
		slog.Error("Error creating chat completion", errors.New("no choices"))
		return errors.New("no choices")
	}
	fmt.Printf("AI> %v\n", response.Choices[0].Message.Content)
	return nil
}

// checkInput checks if text is in the correct format
func (s *testGenService) checkInput(text string) error {
	if !strings.HasPrefix(text, ":testgen") {
		return errors.New("invalid format: text must start with ':testgen'")
	}
	tokens := strings.Split(text, " ")
	if len(tokens) < 2 {
		return errors.New("invalid format: text must contain a file name or a file name and a function name")
	}
	return nil
}

// parseInput parses text
func (s *testGenService) parseInput(text string) (string, string, error) {
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

// SendRequestStream sends request to OpenAI to generate test code in stream
// This expects text to be in the following format:
// :testgen <file> or :testgen <file> <function>
func (s *testGenService) SendRequestStream(ctx context.Context, text string) error {
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
	messageBody := fmt.Sprintf("%s\n\n%s", tesgGenMessageHeader, code)

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
