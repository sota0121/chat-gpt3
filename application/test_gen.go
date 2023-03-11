package application

import (
	"context"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
	"github.com/sota0121/go-ai-chat/internal"
	"golang.org/x/exp/slog"
)

type TestGenService interface {
	SendRequest(ctx context.Context, text string) error
}

func NewTestGenService() TestGenService {
	return &testGenService{}
}

type testGenService struct {
}

var _ TestGenService = (*testGenService)(nil)

const (
	messageHeader = "以下のプログラムについて、テストコードを生成してください。"
)

// SendRequest sends request to OpenAI to generate test code
// This expects text to be in the following format:
// :testgen <file> or :testgen <file> <function>
func (s *testGenService) SendRequest(ctx context.Context, text string) error {
	// Check if text is in the correct format
	if !strings.HasPrefix(text, ":testgen") {
		return errors.New("invalid format: text must start with ':testgen'")
	}
	tokens := strings.Split(text, " ")
	if len(tokens) < 2 {
		return errors.New("invalid format: text must contain a file name or a file name and a function name")
	}

	// Parse text
	fileName := tokens[1]
	funcName := ""
	if len(tokens) > 2 {
		funcName = tokens[2]
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
	messageBody := fmt.Sprintf("%s\n\n%s", messageHeader, code)

	// debug
	fmt.Println(messageBody)
	// debug

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

// extractCode extracts code from file
// This uses AST if a function name is provided
// Otherwise, this reads the whole file
func extractCode(file *os.File, funcName string) (string, error) {
	// Read the whole file
	if funcName == "" {
		return getFileContent(file)
	}

	// Parse file and extract the target function
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, file.Name(), nil, parser.ParseComments)
	if err != nil {
		return "", err
	}
	for _, decl := range f.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			if fn.Name.Name == funcName {
				// Get where the function starts and ends
				start := fset.Position(fn.Pos()).Line
				end := fset.Position(fn.End()).Line

				// Get the content of the file
				content, err := getFileContent(file)
				if err != nil {
					slog.Error("Error getting file content", err)
					return "", err
				}

				// Extract the target function
				lines := strings.Split(content, "\n")
				return strings.Join(lines[start-1:end], "\n"), nil
			}
		}
	}
	err = errors.New("function not found")
	return "", err
}

// getFileContent returns the content of the file
func getFileContent(file *os.File) (string, error) {
	fi, err := file.Stat()
	if err != nil {
		return "", err
	}
	data := make([]byte, fi.Size())
	_, err = file.Read(data)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
