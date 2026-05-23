package cli

import (
	"fmt"

	appconfig "github.com/sorabhlahoti/devmate-ai/internal/config"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage DevMate CLI configuration",
}

var configSetCmd = &cobra.Command{
	Use:     "set [key] [value]",
	Short:   "Set a config value",
	Example: `devmate config set api_url http://localhost:8080`,
	Args:    cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value := args[1]

		if err := appconfig.SetValue("", key, value); err != nil {
			return err
		}

		fmt.Println("Config updated successfully ?")
		fmt.Println(key + " = " + value)

		return nil
	},
}

var configGetCmd = &cobra.Command{
	Use:     "get [key]",
	Short:   "Get a config value",
	Example: `devmate config get api_url`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]

		value, err := appconfig.GetValue("", key)
		if err != nil {
			return err
		}

		fmt.Println(value)

		return nil
	},
}

var configPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Print config file path",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := appconfig.DefaultPath()
		if err != nil {
			return err
		}

		fmt.Println(path)

		return nil
	},
}

func init() {
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configPathCmd)

	rootCmd.AddCommand(configCmd)
}
