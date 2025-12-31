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
	return hash.Sum([]byte(normalized))
}
