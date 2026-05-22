package ai

import (
	"strings"
	"testing"
)

func TestBuildErrorExplanationPrompt(t *testing.T) {
	prompt, err := BuildErrorExplanationPrompt(" permission denied  ")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !strings.Contains(prompt, "permission denied") {
		t.Fatalf("expected prompt to contain error text, got %s", prompt)
	}

	expectedSections := []string{
		"1. What the error means",
		"2. Most likely reason",
		"3. Step-by-step fix",
		"4. How to prevent it next time",
	}

	for _, section := range expectedSections {
		if !strings.Contains(prompt, section) {
			t.Fatalf("expected prompt to contain section %q, got %s", section, prompt)
		}
	}
}

func TestBuildErrorExplanationPromptEmptyText(t *testing.T) {
	_, err := BuildErrorExplanationPrompt("   ")
	if err == nil {
		t.Fatal("expected error for empty error text, got nil")
	}
}
