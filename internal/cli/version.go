package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print DevMate AI version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("DevMate AI version:", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
