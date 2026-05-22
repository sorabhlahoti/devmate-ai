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

    if strings.Contains(cleanPrompt, "Explain the following terminal error") {
        return `[Mock AI Error Explanation]

1. What the error means
The program tried to do something but the operating system rejected it.

2. Most likely reason
This usually happens because of permission, wrong file path, missing file, or restricted access.

3. Step-by-step fix
- Read the exact error message carefully
- Check if the file or folder exists
- Check permission of the file or folder
- Run the command again after fixing permission/path
- If running in Docker, check volume mount and container user

4. How to prevent it next time
Use clear logs, validate file paths, handle errors properly, and document required permissions.`, nil
    }

    answer := fmt.Sprintf("[Mock AI] You asked: %s\nThis is a mock response. Real AI integration will be added later.", cleanPrompt)

    return answer, nil
}
