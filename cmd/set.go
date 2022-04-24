/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"
	"fmt"

	"github.com/mostfunkyduck/vetinari/internal"
	"github.com/99designs/keyring"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Sets a new value in the keyring",
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		ring, err := keyring.Open(keyring.Config{
			ServiceName: viper.GetString(internal.SERVICENAME),
			KeyCtlScope: viper.GetString(internal.KEYCTLSCOPE),
			KeychainName: viper.GetString(internal.KEYCHAINNAME),
		})
		if (err != nil ) {
			fmt.Printf("could not open keyring: %s\n", err)
			os.Exit(1)
		}
	
		err = ring.Set(keyring.Item{
			Key:  args[0],
			Data: []byte(args[1]),
		})

		if err != nil {
			fmt.Printf("could not set key '%s': %s\n", args[0], err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
}
