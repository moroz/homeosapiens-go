package main

import (
	"fmt"
	"log"

	keychain "github.com/EikaGruppen/go-macos-keychain"
)

func main() {
	client := keychain.NewKeychainClient("homeosapiens")
	name, value := "example", "test_value"
	if err := client.Update(name, value); err != nil {
		log.Fatal(err)
	}

	val, err := client.Get(name)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(val)
}
