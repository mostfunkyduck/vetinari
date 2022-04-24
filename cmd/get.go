package cmd

import (
	"os"
	"fmt"

	"github.com/mostfunkyduck/vetinari/internal"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "retrieves a value from the keyring, prints to STDOUT (handle with care)",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ring, err := internal.OpenKeyring()
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not open keyring: %s\n", err)
			os.Exit(1)
		}
		item, err := ring.Get(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not get value for key '%s': %s\n", args[0], err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "%s", item.Data)
		
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
