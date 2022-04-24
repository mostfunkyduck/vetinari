package cmd

import (
	"fmt"
	"os"

	"github.com/mostfunkyduck/vetinari/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)


// rootCmd represents the base command when called without any subcommands
var (
	rootCmd = &cobra.Command{
		Use: internal.NAME,  
		Short: "application for manipulating values in the system's keyring",
	}
	cfgFile	string
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().Bool(internal.UNSAFE, false, "whether or not to print secrets during commands that aren't 'get'")
	rootCmd.PersistentFlags().String(internal.SERVICENAME, internal.NAME, "service name for keyring backends that support this feature")
	rootCmd.PersistentFlags().String(internal.KEYCTLSCOPE, "user", "service name for keyring backends that support this feature")
	rootCmd.PersistentFlags().String(internal.KEYCHAINNAME, internal.NAME, "name of keychain")

	for _, each := range []string{ internal.UNSAFE, internal.SERVICENAME, internal.KEYCTLSCOPE, internal.KEYCHAINNAME } {
		if err := viper.BindPFlag(each, rootCmd.PersistentFlags().Lookup(each)); err != nil {
			fmt.Fprintf(os.Stderr, "error configuring flag '%s', %s\n", each, err)
			os.Exit(1)
		}
	}
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".vetinari")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
