package internal

import (
	"github.com/99designs/keyring"
	"github.com/spf13/viper"
)

func OpenKeyring() (keyring.Keyring, error) {
	return keyring.Open(keyring.Config{
		ServiceName: viper.GetString(SERVICENAME),
		KeyCtlScope: viper.GetString(KEYCTLSCOPE),
		KeychainName: viper.GetString(KEYCHAINNAME),
	})
}
