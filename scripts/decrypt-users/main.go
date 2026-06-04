package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/types"
)

func resolveGitDir() (string, error) {
	start, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dir := start

	for {
		gitDir := filepath.Join(dir, ".git")
		stat, _ := os.Stat(gitDir)
		if stat != nil && stat.IsDir() {
			return dir, nil
		}
		if filepath.Dir(dir) == dir {
			return "", fmt.Errorf("Git directory not found in parent directories of %s", start)
		}
		dir = filepath.Dir(dir)
	}
}

func resolveSopsExecutable() (string, error) {
	return exec.LookPath("sops")
}

func keyFilePresent(gitRoot string) (bool, error) {
	info, err := os.Stat(filepath.Join(gitRoot, "age.txt"))
	if err != nil {
		return false, err
	}

	return !info.IsDir(), nil
}

func decryptSeedFile() (string, error) {
	executable, err := resolveSopsExecutable()
	if err != nil {
		return "", err
	}

	gitRoot, err := resolveGitDir()
	if err != nil {
		return "", err
	}

	keyPresent, err := keyFilePresent(gitRoot)
	if err != nil || !keyPresent {
		return "", err
	}

	seedFile := filepath.Join(gitRoot, "db/seeds/users.csv.enc")

	var out bytes.Buffer

	cmd := exec.Command(executable, "decrypt", seedFile)
	cmd.Dir = gitRoot
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}

	return out.String(), nil
}

func main() {
	db, err := pgxpool.New(context.Background(), config.DatabaseUrl)
	if err != nil {
		log.Fatal(err)
	}

	decrypted, err := decryptSeedFile()
	if err != nil {
		log.Fatal(err)
	}

	tx, err := db.Begin(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback(context.Background())

	reader := csv.NewReader(bytes.NewBufferString(decrypted))

	for {
		row, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}

		name := strings.TrimSpace(row[2])
		givenName, familyName, _ := strings.Cut(name, " ")

		_, err = services.NewUserService(tx).CreateUser(context.Background(), &types.SeedUserParams{
			EmailConfirmed: true,
			GivenName:      givenName,
			FamilyName:     familyName,
			Email:          row[1],
			Country:        "PL",
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	tx.Commit(context.Background())
}
