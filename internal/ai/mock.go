package ai

import (
    "context"
    "fmt"
    "strings"
)

type MockProvider struct{}

func NewMockProvider() Provider {
    return &MockProvider{}
}

func (m *MockProvider) Ask(ctx context.Context, prompt string) (string, error) {
    cleanPrompt := strings.TrimSpace(prompt)

    if cleanPrompt == "" {
        return "", fmt.Errorf("prompt cannot be empty")
    }

    answer := fmt.Sprintf("[Mock AI] You asked: %s\nThis is a mock response. Real AI integration will be added later.", cleanPrompt)

    return answer, nil
}
