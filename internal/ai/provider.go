package ai

import "context"

type Provider interface {
    Ask(ctx context.Context, prompt string) (string, error)
}
