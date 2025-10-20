package llm

import "context"

type Provider interface {
	GenerateCode(ctx context.Context, prompt string) (string, error)
}