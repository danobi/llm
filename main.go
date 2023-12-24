package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func main() {
	// Create client given API key
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("API_KEY")))
	if err != nil {
		log.Fatal("Failed to create genai client: %v", err)
	}
	defer client.Close()

	// Get input from stdin
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal("Failed to read stdin: %v", err)
	}

	// Query model
	model := client.GenerativeModel("models/gemini-pro")
	resp, err := model.GenerateContent(ctx, genai.Text(input))
	if err != nil {
		log.Fatal("Failed to generate response: %v", err)
	}

	// Parse response
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) != 1 {
		fmt.Fprintf(os.Stderr, "No response received\n")
		return
	}
	fmt.Println(resp.Candidates[0].Content.Parts[0])
}
