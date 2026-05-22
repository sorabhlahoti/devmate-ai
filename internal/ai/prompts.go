package ai

import (
	"fmt"
	"strings"
)

func BuildErrorExplanationPrompt(errorText string) (string, error) {
	cleanText := strings.TrimSpace(errorText)

	if cleanText == "" {
		return "", fmt.Errorf("error text cannot be empty")
	}

	prompt := fmt.Sprintf(`You are DevMate AI, a developer debugging assistant.

Explain the following terminal error in simple words.

Error:
%s

Return answer in this format:
1. What the error means
2. Most likely reason
3. Step-by-step fix
4. How to prevent it next time
`, cleanText)

	return prompt, nil
}
