package main

import (
	"fmt"
	"os"

	"github.com/moroz/homeosapiens-go/bin/cli/client"
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

func GetExistingToken() (string, error) {
	keyringClient, err := client.NewDefaultKeyringClient()
	if err != nil {
		return "", err
	}

	return keyringClient.FetchAPIToken(config.ApiTokenKeyringKey)
}

func SignIn() (string, error) {
	keyringClient, err := client.NewDefaultKeyringClient()
	if err != nil {
		return "", err
	}

	var newToken string
	for {
		newToken, err = readApiTokenFromStdin("Please paste your API token (will not show in terminal): ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading token: %v. Please try again.", err)
		}
		if newToken != "" {
			break
		}
	}

	if err := keyringClient.SetAPIToken(config.ApiTokenKeyringKey, newToken); err != nil {
		return "", err
	}

	return newToken, nil
}

func EnsureSignedIn() (string, error) {
	token, err := GetExistingToken()
	if err != nil {
		return "", err
	}
	if token == "" {
		fmt.Fprintf(os.Stderr, "You are not signed in. Press Enter to open the browser.")

	}
}

func main() {

}
