package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
	"github.com/sota0121/go-ai-chat/application"
	"github.com/sota0121/go-ai-chat/internal"
	"golang.org/x/exp/slog"
	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Commands map[string]Command `yaml:"commands"`
}

type Command struct {
	SystemMessages []string `yaml:"systemMessages"`
	UserMessages   []string `yaml:"userMessages"`
}

const (
	envFileName         = ".env"
	configFileName      = "config.yml"
	openAiApiKeyEnvName = "OPENAI_API_KEY"
	keyCommandsChat     = "chat"
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

	// Load config file
	cfg, err := loadConfig(filepath.Join(filepath.Dir(exec), configFileName))
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
	ctx = internal.SetOpenAIClientToContext(ctx, openaiClient)
	app := NewApp(ctx, cfg)

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

// loadConfig loads config file
func loadConfig(fileName string) (*Config, error) {
	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		slog.Error("yamlFile.Get err   #%v ", err)
		return nil, err
	}

	var config Config

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		slog.Error("Unmarshal: %v", err)
		return nil, err
	}
	return &config, nil
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
	config         *Config
	CommandService application.CommandService
	ChatService    application.ChatService
	FindBugService application.FindBugService
	TestGenService application.TestGenService
}

func NewApp(ctx context.Context, cfg *Config) *App {
	systemMessages := cfg.Commands[keyCommandsChat].SystemMessages
	userMessages := cfg.Commands[keyCommandsChat].UserMessages

	return &App{
		ctx:            ctx,
		config:         cfg,
		CommandService: application.NewCommandService(),
		ChatService:    application.NewChatService(systemMessages, userMessages),
		FindBugService: application.NewFindBugService(),
		TestGenService: application.NewTestGenService(),
	}
}

// Execute executes application
func (a *App) Execute() error {
	fmt.Println("Please input text (or ':quit' to exit): ")

	// Process loop
	for {
		fmt.Printf("chat> ")
		s := bufio.NewScanner(os.Stdin)
		s.Scan()

		// Parse command
		if strings.HasPrefix(s.Text(), ":") {
			ct := a.CommandService.ParseCommand(s.Text())
			switch ct {
			case application.ShowHelp:
				a.CommandService.ShowHelp()
			case application.ShowVersion:
				a.CommandService.ShowVersion()
			case application.Quit:
				fmt.Println("Bye!")
				return nil
			case application.TestGen:
				err := a.TestGenService.SendRequestStream(a.ctx, s.Text())
				if err != nil {
					slog.Error("Error TestGenService.SendRequestStream", err)
					break
				}
			case application.FindBugs:
				err := a.FindBugService.SendRequestStream(a.ctx, s.Text())
				if err != nil {
					slog.Error("Error FindBugService.SendRequestStream", err)
					break
				}
			}
			continue
		}
		a.ChatService.SendTextStream(a.ctx, s.Text())
		// a.ChatService.SendText(a.ctx, s.Text())
	}
}
