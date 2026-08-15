package main

import (
	"fmt"
	"log"
	"os"

	"github.com/moroz/homeosapiens-go/bin/cli/auth"
	"github.com/moroz/homeosapiens-go/bin/cli/config"
	"golang.org/x/term"
)

func readApiTokenFromStdin(prompt string) (string, error) {
	fmt.Print(prompt)

	bytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func main() {
	client, err := auth.NewClient(config.KeychainName, config.ServiceName)
	if err != nil {
		log.Fatal(err)
	}

	token, err := client.FetchAPIToken(config.ApiTokenKeyringKey)
	if err != nil {
		log.Fatal(err)
	}

	if token != "" {
		fmt.Println(token)
		return
	}

	var newToken string
	for {
		newToken, err = readApiTokenFromStdin("Please paste your API token (will not show in terminal): ")
		if err != nil {
			log.Fatal(err)
		}
		if newToken != "" {
			break
		}
	}
	if err := client.SetAPIToken(config.ApiTokenKeyringKey, newToken); err != nil {
		log.Fatal(err)
	}
}
