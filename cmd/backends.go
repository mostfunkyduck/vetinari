package cmd

import (
	"fmt"

	"github.com/99designs/keyring"
	"github.com/spf13/cobra"
)

// backendsCmd represents the backends command
var backendsCmd = &cobra.Command{
	Use:   "backends",
	Short: "Lists all available keychain backends in usage priority order",
	Run: func(cmd *cobra.Command, args []string) {
		for _, backend := range keyring.AvailableBackends() {
			fmt.Println(backend)
		}
	},
}

func init() {
	rootCmd.AddCommand(backendsCmd)
}
