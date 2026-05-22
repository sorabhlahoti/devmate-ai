package cli

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sorabhlahoti/devmate-ai/internal/ai"

	"github.com/spf13/cobra"
)

var errorFilePath string

var explainErrorCmd = &cobra.Command{
	Use:     "explain-error",
	Short:   "Explain a terminal or application error from a file",
	Example: "devmate explain-error --file error.log",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if strings.TrimSpace(errorFilePath) == "" {
			return fmt.Errorf("please provide an error file using --file")
		}

		content, err := os.ReadFile(errorFilePath)
		if err != nil {
			return fmt.Errorf("failed to read error file: %w", err)
		}

		prompt, err := ai.BuildErrorExplanationPrompt(string(content))
		if err != nil {
			return err
		}

		provider := ai.NewMockProvider()

		answer, err := provider.Ask(context.Background(), prompt)
		if err != nil {
			return err
		}

		fmt.Println("Error File:")
		fmt.Println(errorFilePath)
		fmt.Println()
		fmt.Println("Explanation:")
		fmt.Println(answer)

		return nil
	},
}

func init() {
	explainErrorCmd.Flags().StringVarP(&errorFilePath, "file", "f", "", "Path to error log file")
	rootCmd.AddCommand(explainErrorCmd)
}
