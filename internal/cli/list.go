package cli

import (
	"context"
	"fmt"

	"github.com/sorabhlahoti/devmate-ai/internal/storage"

	"github.com/spf13/cobra"
)

var listLimit int

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List saved developer notes",
	Example: `devmate list --limit 10`,
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := storage.NewSQLiteStore("")
		if err != nil {
			return err
		}
		defer store.Close()

		notes, err := store.ListNotes(context.Background(), listLimit)
		if err != nil {
			return err
		}

		if len(notes) == 0 {
			fmt.Println("No notes found.")
			return nil
		}

		for _, note := range notes {
			fmt.Printf("[%d] %s\n", note.ID, note.Title)
			fmt.Println(note.Body)
			fmt.Println("Created:", note.CreatedAt.Format("2006-01-02 15:04:05 UTC"))

			if note.RemoteID != nil && note.SyncedAt != nil {
				fmt.Printf("Sync: synced as remote note %d at %s\n", *note.RemoteID, note.SyncedAt.Format("2006-01-02 15:04:05 UTC"))
			} else {
				fmt.Println("Sync: not synced")
			}

			fmt.Println("---")
		}

		return nil
	},
}

func init() {
	listCmd.Flags().IntVarP(&listLimit, "limit", "l", 20, "Maximum number of notes to list")

	rootCmd.AddCommand(listCmd)
}
