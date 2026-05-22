package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/sorabhlahoti/devmate-ai/internal/ai"

	"github.com/spf13/cobra"
)

var askCmd = &cobra.Command{
	Use:     "ask [question]",
	Short:   "Ask a technical question from the terminal",
	Example: "devmate ask \"explain goroutine in simple words\"",
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		question := strings.Join(args, " ")

		provider := ai.NewMockProvider()

		answer, err := provider.Ask(context.Background(), question)
		if err != nil {
			return err
		}

		fmt.Println("Question:")
		fmt.Println(question)
		fmt.Println()
		fmt.Println("Answer:")
		fmt.Println(answer)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(askCmd)
}
