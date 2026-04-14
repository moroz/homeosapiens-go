package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"log"
	"strings"

	"github.com/bincyber/go-sqlcrypter"
	"github.com/moroz/homeosapiens-go/config"
)

func init() {
	crypterer, err := NewEncryptionProvider(config.DatabaseEncryptionKey, nil)
	if err != nil {
		log.Fatalf("Failed to initialize database encryption provider: %s", err)
	}
	sqlcrypter.Init(crypterer)
}

func HashEmail(email string) []byte {
	normalized := strings.TrimSpace(strings.ToLower(email))
	hash := hmac.New(sha256.New, config.DatabaseHMACKey)
	hash.Write([]byte(normalized))
	return hash.Sum(nil)
}

func HashUserToken(token []byte) []byte {
	hash := hmac.New(sha256.New, config.DatabaseHMACKey)
	hash.Write(token)
	return hash.Sum(nil)
}
