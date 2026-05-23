package ai

import (
	"context"
	"strings"
	"testing"
)

func TestMockProviderAsk(t *testing.T) {
	provider := NewMockProvider()

	answer, err := provider.Ask(context.Background(), "explain goroutine")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !strings.Contains(answer, "explain goroutine") {
		t.Fatalf("expected answer to contain prompt, got %s", answer)
	}
}

func TestMockProviderAskEmptyPrompt(t *testing.T) {
	provider := NewMockProvider()

	_, err := provider.Ask(context.Background(), "   ")
	if err == nil {
		t.Fatal("expected error for empty prompt, got nil")
	}
}
