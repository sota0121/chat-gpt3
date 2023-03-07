package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
	"github.com/sota0121/chat-gpt3/example"
)

func main() {
	// Load .env file
	fileName := ".env"
	err := godotenv.Load(fileName)
	if err != nil {
		log.Fatal("Error loading .env file")
		os.Exit(1)
	}

	// Get OpenAI API Key
	openaiApiKey := os.Getenv("OPENAI_API_KEY")
	hiddenApiKey := openaiApiKey[:4] + strings.Repeat("*", len(openaiApiKey)-4)
	fmt.Println("OpenAI API Key: ", hiddenApiKey)

	// Create OpenAI client only once
	openaiClient := openai.NewClient(openaiApiKey)

	// Invoke stream example
	// example.ExampleChatStream(openaiClient)

	// Invoke unary example
	fmt.Println("Start chat with GPT-3 (unary example)...")
	s := bufio.NewScanner(os.Stdin)
	s.Scan()
	example.ExampleChat(openaiClient, s.Text())

}
