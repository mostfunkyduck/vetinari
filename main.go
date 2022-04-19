package main

import (
	"flag"
	"fmt"
	"github.com/99designs/keyring"
	"os"
	"regexp"
)

type Command string

const (
	GET          Command = "get"
	SET          Command = "set"
	LIST         Command = "list"
	LISTBACKENDS Command = "backends"
)

type Config struct {
	// Which command was specified
	Command Command
	// configuration in case of a get command
	GetConfig GetConfig
	// configuration for keyring library
	KeyringConfig keyring.Config
	PrintSecrets  bool
	SetConfig     SetConfig
}

type SetConfig struct {
	Key   string
	Value string
}

type GetConfig struct {
	Key string
}

type GlobalFlags struct {
	ServiceName  *string
	KeyCtlScope  *string
	KeychainName *string
	PrintSecrets *bool
}

func generateGlobalFlags(executableName string, fs *flag.FlagSet) (GlobalFlags, error) {
	globalFlags := GlobalFlags{
		ServiceName:  fs.String("service", executableName, "set service for keyrings that support it"),
		KeyCtlScope:  fs.String("keyctl-scope", "user", "set scope for keyctl, if it is in use. can be 'session', 'user', 'process', or 'thread'"),
		KeychainName: fs.String("keychain-name", executableName, "name of the keychain used on macos"),
		PrintSecrets: fs.Bool("print-secrets", false, "whether or not to print secrets during the operation of commands besides 'get', for debugging"),
	}
	return globalFlags, nil
}

func parseFlags() (Config, error) {
	var command Command
	executableName, err := os.Executable()
	executableName = regexp.MustCompile(`.*\/`).ReplaceAllString(executableName, "")
	if err != nil {
		return Config{}, fmt.Errorf("could not get name of executable for usage in flag default: %s", err)
	}
	setSet := flag.NewFlagSet("set", flag.ExitOnError)
	getSet := flag.NewFlagSet("get", flag.ExitOnError)
	listSet := flag.NewFlagSet("list", flag.ExitOnError)

	var setConfig SetConfig
	var getConfig GetConfig
	var globalFlags GlobalFlags

	switch os.Args[1] {
	case "set":
		command = SET
		key := setSet.String("key", "", "key in keyring to set value for, required")
		value := setSet.String("value", "", "value to set in keyring")

		gf, err := generateGlobalFlags(executableName, setSet)
		if err != nil {
			return Config{}, fmt.Errorf("could not generate global flags for 'set': %s", err)
		}
		globalFlags = gf
		setSet.Usage = func() {
			fmt.Fprintf(os.Stderr, "sets a value in the keychain\n")
			fmt.Fprintf(os.Stderr, "\n")
			fmt.Fprintf(os.Stderr, "respects all arguments set at the main command")
			fmt.Fprintf(os.Stderr, "\n")
			setSet.PrintDefaults()
		}
		if err := setSet.Parse(os.Args[2:]); err != nil {
			return Config{}, fmt.Errorf("could not parse arguments for 'set': %s", err)
		}
		setConfig = SetConfig{
			Key:   *key,
			Value: *value,
		}
		if setConfig.Key == "" {
			return Config{}, fmt.Errorf("'set' was specified, but no -key argument was passed in")
		}
	case "backends":
		command = LISTBACKENDS
	case "list":
		command = LIST
		gf, err := generateGlobalFlags(executableName, listSet)
		if err != nil {
			return Config{}, fmt.Errorf("could not generate global flags for 'list': %s", err)
		}
		globalFlags = gf

		if err := listSet.Parse(os.Args[2:]); err != nil {
			return Config{}, fmt.Errorf("could not parse arguments for 'list': %s", err)
		}
	case "get":
		command = GET
		getSet.Usage = func() {
			fmt.Fprintf(os.Stderr, "gets a value in the keychain\n")
			fmt.Fprintf(os.Stderr, "\n")
			fmt.Fprintf(os.Stderr, "respects all arguments set at the main command")
			fmt.Fprintf(os.Stderr, "\n")
			setSet.PrintDefaults()
		}
		key := getSet.String("key", "", "key in keyring to retrieve value for, required")
		if err := getSet.Parse(os.Args[2:]); err != nil {
			return Config{}, fmt.Errorf("could not parse arguments for 'get': %s", err)
		}
		getConfig = GetConfig{
			Key: *key,
		}

		if getConfig.Key == "" {
			return Config{}, fmt.Errorf("'set' was specified, but no -key argument was passed in")
		}
	}
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s will interact with various keyrings in order to securely manage secrets\n", executableName)
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "the direct arguments are:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "in addition, the following subcommands are supported:\n")
		fmt.Fprintf(os.Stderr, "set\t\t\t- sets a value in the keyring\n")
		fmt.Fprintf(os.Stderr, "get\t\t\t- retrieves a value in the keyring\n")
		fmt.Fprintf(os.Stderr, "backends\t\t- lists all the backends for the keyring ordered by what the utility prefers\n")
	}

	flag.Parse()

	config := Config{
		Command:      command,
		SetConfig:    setConfig,
		GetConfig:    getConfig,
		PrintSecrets: *globalFlags.PrintSecrets,
		KeyringConfig: keyring.Config{
			ServiceName:  *globalFlags.ServiceName,
			KeyCtlScope:  *globalFlags.KeyCtlScope,
			KeychainName: *globalFlags.KeychainName,
		},
	}
	return config, nil
}

func listBackends() []string {
	var returnVal []string
	for _, backend := range keyring.AvailableBackends() {
		returnVal = append(returnVal, string(backend))
	}
	return returnVal
}

func main() {

	if len(os.Args) <= 1 {
		fmt.Fprintf(os.Stderr, "no subcommand specified\n")
		os.Exit(1)
	}

	config, err := parseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse flags: %s\n", err)
		os.Exit(1)
	}

	if config.Command == LISTBACKENDS {
		for _, backend := range listBackends() {
			fmt.Printf("\t%s\n", backend)
		}
		os.Exit(0)
	}

	ring, err := keyring.Open(config.KeyringConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not open keyring: %s\n", err)
		os.Exit(1)
	}

	switch config.Command {
	case LIST:
		keys, err := ring.Keys()
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not retrieve keys from keyring: %s\n", err)
			os.Exit(1)
		}

		for _, key := range keys {
			value := "****"
			if config.PrintSecrets {
				item, err := ring.Get(key)
				if err != nil {
					fmt.Fprintf(os.Stderr, "could not retrieve key '%s' from keyring: %s\n", key, err)
					os.Exit(1)
				}
				value = string(item.Data)
			}
			fmt.Printf("%s: %s\n", key, value)
		}
	case SET:
		err = ring.Set(keyring.Item{
			Key:  config.SetConfig.Key,
			Data: []byte(config.SetConfig.Value),
		})

		if err != nil {
			fmt.Printf("could not set key '%s': %s\n", config.SetConfig.Key, err)
			os.Exit(1)
		}

	case GET:
		i, err := ring.Get(config.GetConfig.Key)
		if err != nil {
			fmt.Printf("could not get value for key '%s': %s\n", config.GetConfig.Key, err)
			os.Exit(1)
		}
		fmt.Printf("%s", i.Data)
	}
}
