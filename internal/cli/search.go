package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/sorabhlahoti/devmate-ai/internal/storage"

	"github.com/spf13/cobra"
)

var searchLimit int

var searchCmd = &cobra.Command{
	Use:     "search [text]",
	Short:   "Search saved developer notes",
	Example: `devmate search "docker permission"`,
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		searchText := strings.Join(args, " ")

		store, err := storage.NewSQLiteStore("")
		if err != nil {
			return err
		}
		defer store.Close()

		notes, err := store.SearchNotes(context.Background(), searchText, searchLimit)
		if err != nil {
			return err
		}

		if len(notes) == 0 {
			fmt.Println("No matching notes found.")
			return nil
		}

		fmt.Println("Search results:")
		fmt.Println()

		for _, note := range notes {
			fmt.Printf("[%d] %s\n", note.ID, note.Title)
			fmt.Println(note.Body)
			fmt.Println("Created:", note.CreatedAt.Format("2006-01-02 15:04:05 UTC"))
			fmt.Println("---")
		}

		return nil
	},
}

func init() {
	searchCmd.Flags().IntVarP(&searchLimit, "limit", "l", 20, "Maximum number of notes to search")

	rootCmd.AddCommand(searchCmd)
}
