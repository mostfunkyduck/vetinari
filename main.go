package main
import (
	"fmt"
	"os"
	"github.com/99designs/keyring"
)

func main() {
	ring, err := keyring.Open(keyring.Config{
		ServiceName: "example",
		KeyCtlScope: "user",
	})

	if (err != nil) {
		fmt.Printf("could not open keyring: %s\n", err)
		os.Exit(1)
	}
	err = ring.Set(keyring.Item{
		Key: "foo",
		Data: []byte("secret-bar"),
	})

	if (err != nil) {
		fmt.Printf("could not set : %s\n", err)
		os.Exit(1)
	}

	i, err := ring.Get("foo")
	if (err != nil) {
		fmt.Printf("could not get: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("%s", i.Data) 
}
