package main

import (
	"context"
	_ "embed"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

//go:embed main.go
var src string

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

// Queries model and prints out result to stdout if successful
func ask(in []genai.Part) error {
	ctx := context.Background()

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

	// Query model
	model := client.GenerativeModel("models/gemini-pro")
	resp, err := model.GenerateContent(ctx, in...)
	if err != nil {
		return fmt.Errorf("Failed to generate response: %v", err)
	}

	// Parse response
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) != 1 {
		return fmt.Errorf("no response received")
	}

	fmt.Println(resp.Candidates[0].Content.Parts[0])
	return nil
}

// Prints help message tou stderr
func help() {
	// Boring, non-innovate help message
	if innovate := os.Getenv("INNOVATE"); innovate == "" {
		fmt.Fprintf(os.Stderr, "llm [-][context]..\n")
		return
	}

	// Now we're talking - it's time to innovate
	prompt := `
		Given the following go program that compiles to a binary titled
		'llm', print out an appropriate --help message.  Your response
		should be human readable text suitable to immediately print to
		terminal.  Do not use any code blocks or backticks.  Only show
		short form help text which at minimum explains the positional
		parameters. Do not show anything like a man page but provide
		at least a sentence or two.
	`
	parts := []genai.Part{
		genai.Text(prompt),
		genai.Text(src),
	}
	if err := ask(parts); err != nil {
		log.Fatalf("Failed to innovate: %v", err)
	}
}

func main() {
	// Handle help message
	for _, arg := range os.Args {
		if arg == "-h" || arg == "--help" {
			help()
			return
		}
	}

	// Acquire input from user
	in, err := input()
	if err != nil {
		log.Fatalf("Failed to get input: %v", err)
	}

	// Ask gemini
	if err := ask(in); err != nil {
		log.Fatal(err)
	}
}
