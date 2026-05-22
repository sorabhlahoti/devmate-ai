package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "0.1.0"

var rootCmd = &cobra.Command{
	Use:   "devmate",
	Short: "DevMate AI is an AI-powered CLI assistant for developers",
	Long: "DevMate AI helps developers ask questions, explain errors, " +
		"save useful fixes, and generate developer runbooks from the terminal.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to DevMate AI ??")
		fmt.Println("Run 'devmate --help' to see available commands.")
	},
}

func Execute() error {
	return rootCmd.Execute()
}
