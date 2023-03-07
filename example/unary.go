package example

import (
	"context"
	"fmt"
	"log"

	"github.com/sashabaranov/go-openai"
)

// ExampleChatStream is an example of how to use the stream API
func ExampleChat(openaiClient *openai.Client, input string) {
	fmt.Println("in: ", input)

	// Make a request
	ctx := context.Background()

	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: input,
			},
		},
	}

	response, err := openaiClient.CreateChatCompletion(ctx, req)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Printf("Response: %v\n", response)
}
