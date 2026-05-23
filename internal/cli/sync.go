package cli

import (
	"context"
	"fmt"

	"github.com/sorabhlahoti/devmate-ai/internal/apiclient"
	appconfig "github.com/sorabhlahoti/devmate-ai/internal/config"
	"github.com/sorabhlahoti/devmate-ai/internal/storage"

	"github.com/spf13/cobra"
)

var syncLimit int

var syncCmd = &cobra.Command{
	Use:     "sync",
	Short:   "Sync unsynced local notes with backend API",
	Example: `devmate sync --limit 20`,
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		cfg, err := appconfig.Load("")
		if err != nil {
			return err
		}

		client, err := apiclient.New(cfg.APIURL)
		if err != nil {
			return err
		}

		store, err := storage.NewSQLiteStore("")
		if err != nil {
			return err
		}
		defer store.Close()

		notes, err := store.ListUnsyncedNotes(ctx, syncLimit)
		if err != nil {
			return err
		}

		if len(notes) == 0 {
			fmt.Println("No unsynced notes found ?")
			return nil
		}

		fmt.Println("Syncing notes to:", cfg.APIURL)
		fmt.Println()

		successCount := 0

		for _, note := range notes {
			remoteNote, err := client.CreateNote(ctx, note.Title, note.Body)
			if err != nil {
				fmt.Printf("Failed to sync local note %d: %v\n", note.ID, err)
				continue
			}

			if err := store.MarkNoteSynced(ctx, note.ID, remoteNote.ID); err != nil {
				return err
			}

			successCount++

			fmt.Printf("Synced local note %d -> remote note %d ?\n", note.ID, remoteNote.ID)
		}

		fmt.Println()
		fmt.Printf("Sync complete: %d/%d notes synced\n", successCount, len(notes))

		if successCount != len(notes) {
			return fmt.Errorf("some notes failed to sync")
		}

		return nil
	},
}

func init() {
	syncCmd.Flags().IntVarP(&syncLimit, "limit", "l", 20, "Maximum number of notes to sync")

	rootCmd.AddCommand(syncCmd)
}
