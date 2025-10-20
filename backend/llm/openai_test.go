package llm

import (
	"context"
	"os"
	"testing"
	"time"
)

// TestOpenAIClient_GenerateCode tests the OpenAI client
// Run with: OPENAI_API_KEY=sk-... go test -v ./llm
func TestOpenAIClient_GenerateCode(t *testing.T) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		t.Skip("OPENAI_API_KEY not set, skipping integration test")
	}

	client := NewOpenAIClient(apiKey, "gpt-4")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	prompt := "Create a simple hello world page with a button"

	result, err := client.GenerateCode(ctx, prompt)
	if err != nil {
		t.Fatalf("GenerateCode failed: %v", err)
	}

	if result == "" {
		t.Fatal("Expected non-empty result")
	}

	t.Logf("Generated code length: %d characters", len(result))
	t.Logf("Result preview: %s...", result[:min(200, len(result))])
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}