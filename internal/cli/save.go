package cli

import (
	"context"
	"fmt"

	"github.com/sorabhlahoti/devmate-ai/internal/storage"

	"github.com/spf13/cobra"
)

var saveTitle string
var saveBody string

var saveCmd = &cobra.Command{
	Use:     "save",
	Short:   "Save a developer note or fix locally",
	Example: `devmate save --title "Docker permission fix" --body "Check mounted volume permission"`,
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := storage.NewSQLiteStore("")
		if err != nil {
			return err
		}
		defer store.Close()

		note, err := store.SaveNote(context.Background(), saveTitle, saveBody)
		if err != nil {
			return err
		}

		fmt.Println("Note saved successfully ?")
		fmt.Println("ID:", note.ID)
		fmt.Println("Title:", note.Title)

		return nil
	},
}

func init() {
	saveCmd.Flags().StringVar(&saveTitle, "title", "", "Note title")
	saveCmd.Flags().StringVar(&saveBody, "body", "", "Note body")

	rootCmd.AddCommand(saveCmd)
}
