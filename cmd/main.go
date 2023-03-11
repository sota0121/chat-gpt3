package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
	"github.com/sota0121/go-ai-chat/application"
	"github.com/sota0121/go-ai-chat/example"
	"golang.org/x/exp/slog"
)

const (
	envFileName         = ".env"
	openAiApiKeyEnvName = "OPENAI_API_KEY"
)

func main() {
	// Load .env file
	exec, err := os.Executable()
	if err != nil {
		slog.Error("Error getting executable path", err)
		os.Exit(1)
	}
	err = loadEnv(filepath.Join(filepath.Dir(exec), envFileName))
	if err != nil {
		os.Exit(1)
	}

	// Create OpenAI client
	openaiClient, err := createOpenAIClient()
	if err != nil {
		os.Exit(1)
	}

	// Create application
	ctx := context.Background()
	ctx = SetOpenAIClientToContext(ctx, openaiClient)
	app := NewApp(ctx)

	// Start chat application
	app.Execute()
}

// loadEnv loads .env file
func loadEnv(fileName string) error {
	err := godotenv.Load(fileName)
	if err != nil {
		// .env file is optional so ignore error
		if os.IsNotExist(err) {
			return nil
		}
		slog.Info("Error loading .env file")
		return err
	}
	return nil
}

// createOpenAIClient creates OpenAI client
func createOpenAIClient() (*openai.Client, error) {
	// Get OpenAI API Key
	openaiApiKey := os.Getenv(openAiApiKeyEnvName)
	if openaiApiKey == "" {
		slog.Info("OpenAI API Key is not set")
		return nil, errors.New("OpenAI API Key is not set")
	}

	// Display OpenAI API Key
	hiddenApiKey := openaiApiKey[:4] + strings.Repeat("*", len(openaiApiKey)-4)
	fmt.Println("OpenAI API Key: ", hiddenApiKey)

	// Create OpenAI client only once
	openaiClient := openai.NewClient(openaiApiKey)
	return openaiClient, nil
}

type App struct {
	ctx            context.Context
	ChatService    application.ChatService
	FindBugService application.FindBugService
	TestGenService application.TestGenService
}

func NewApp(ctx context.Context) *App {
	return &App{
		ctx:            ctx,
		ChatService:    application.NewChatService(),
		FindBugService: application.NewFindBugService(),
		TestGenService: application.NewTestGenService(),
	}
}

// Execute executes application
func (a *App) Execute() error {
	fmt.Println("Please input text (or ':quit' to exit): ")
	openaiClient := GetOpenAIClientFromContext(a.ctx)

	// Process loop
	for {
		fmt.Println("chat> ")
		s := bufio.NewScanner(os.Stdin)
		s.Scan()
		if s.Text() == ":quit" {
			break
		}
		example.ExampleChat(openaiClient, s.Text())
	}
	return nil
}
