package cmd

import (
	"os"
	"fmt"

	"github.com/mostfunkyduck/vetinari/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "lists all items on the keyring to STDOUT. If '--unsafe' is specified, lists values as well",
	Run: func(cmd *cobra.Command, args []string) {
		ring, err := internal.OpenKeyring()
		if (err != nil) {
			fmt.Fprintf(os.Stderr, "error opening keyring: %s\n", err)	
		}
		keys, err := ring.Keys()
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not retrieve keys from keyring: %s\n", err)
			os.Exit(1)
		}

		for _, key := range keys {
			value := "****"
			// retrieve the secret if in 'unsafe' mode
			if viper.GetBool(internal.UNSAFE) {
				item, err := ring.Get(key)
				if err != nil {
					fmt.Fprintf(os.Stderr, "could not retrieve key '%s' from keyring: %s\n", key, err)
					os.Exit(1)
				}
				value = string(item.Data)
			}
			fmt.Printf("%s: %s\n", key, value)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
