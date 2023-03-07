package example

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/sashabaranov/go-openai"
)

// ExampleChatStream is an example of how to use the stream API
func ExampleChatStream(openaiClient *openai.Client) {
	// Set up OpenAI client and prepare the stream
	ctx := context.Background()

	req := openai.CompletionRequest{
		Model:     openai.GPT3Dot5Turbo,
		MaxTokens: 5,
		Prompt:    "This is a test",
		Stream:    true,
	}
	stream, err := openaiClient.CreateCompletionStream(ctx, req)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer stream.Close()

	// Start to receive the stream
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("Stream finished")
			return
		}

		if err != nil {
			fmt.Printf("Stream error: %v\n", err)
			return
		}

		fmt.Printf("Stream response: %v\n", response)
	}
}
