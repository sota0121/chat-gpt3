package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	fileName := ".env"
	err := godotenv.Load(fileName)
	if err != nil {
		log.Fatal("Error loading .env file")
		os.Exit(1)
	}

	openaiApiKey := os.Getenv("OPENAI_API_KEY")
	fmt.Println("OpenAI API Key: ", openaiApiKey)
}
