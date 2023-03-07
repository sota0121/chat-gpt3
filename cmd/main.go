package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
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
	fmt.Println("OpenAI API Key: ", openaiApiKey)

	// Invoke stream example
	// example.ExampleChatStream(openaiApiKey)

	// Invoke unary example
	example.ExampleChat(openaiApiKey)

}
