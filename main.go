package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// Get API key from well known locations
func key() (string, error) {
	if env := os.Getenv("API_KEY"); env != "" {
		return env, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home dir: %w", err)
	}

	fpath := filepath.Join(home, ".config/llm/key")
	f, err := os.Open(fpath)
	if err != nil {
		return "", fmt.Errorf("failed to open key file: %w", err)
	}
	defer f.Close()

	key, err := io.ReadAll(f)
	if err != nil {
		return "", fmt.Errorf("failed to read key file: %w", err)
	}

	return string(key), nil
}

func main() {
	ctx := context.Background()

	// Get API key
	k, err := key()
	if err != nil {
		log.Fatal("Failed to get API key: %v", err)
	}

	// Create client given API key
	client, err := genai.NewClient(ctx, option.WithAPIKey(k))
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
