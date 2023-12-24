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

// Get all input from stdin
func stdin() (string, error) {
	in, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", fmt.Errorf("failed to read stdin: %w", err)
	}

	return string(in), nil
}

// Get input from user
func input() ([]genai.Part, error) {
	// Handle all positional parameters
	parts := make([]genai.Part, 0, 1)
	nargs := len(os.Args)
	for i := 1; i < nargs; i++ {
		var part string
		arg := os.Args[i]

		if arg == "-" {
			in, err := stdin()
			if err != nil {
				return nil, err
			}
			part = in
		} else {
			part = os.Args[i]
		}

		parts = append(parts, genai.Text(part))
	}

	// No params were passed, so we have to read from stdin
	if len(parts) == 0 {
		fmt.Fprintf(os.Stderr, "Reading from stdin...\n")
		fmt.Fprintf(os.Stderr, "^C to cancel, ^D to send\n")

		in, err := stdin()
		if err != nil {
			return nil, err
		}

		parts = append(parts, genai.Text(in))
	}

	return parts, nil
}

func main() {
	ctx := context.Background()

	// Handle help message
	for _, arg := range os.Args {
		if arg == "-h" || arg == "--help" {
			fmt.Fprintf(os.Stderr, "llm [-][context]..\n")
			return
		}
	}

	// Get API key
	k, err := key()
	if err != nil {
		log.Fatalf("Failed to get API key: %v", err)
	}

	// Create client given API key
	client, err := genai.NewClient(ctx, option.WithAPIKey(k))
	if err != nil {
		log.Fatalf("Failed to create genai client: %v", err)
	}
	defer client.Close()

	// Acquire input from user
	in, err := input()
	if err != nil {
		log.Fatalf("Failed to get input: %v", err)
	}

	// Query model
	model := client.GenerativeModel("models/gemini-pro")
	resp, err := model.GenerateContent(ctx, in...)
	if err != nil {
		log.Fatalf("Failed to generate response: %v", err)
	}

	// Parse response
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) != 1 {
		fmt.Fprintf(os.Stderr, "No response received\n")
		return
	}
	fmt.Println(resp.Candidates[0].Content.Parts[0])
}
