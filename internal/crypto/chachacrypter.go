package crypto

import (
	"bytes"
	"crypto/cipher"
	"fmt"
	"io"

	"github.com/bincyber/go-sqlcrypter"
	"golang.org/x/crypto/chacha20"
	"golang.org/x/crypto/chacha20poly1305"
)

// ChachaCrypter implements the sqlcrypter.Crypterer interface. The implementation in this file is mostly a blatant copy of https://github.com/bincyber/go-sqlcrypter/blob/master/providers/aesgcm/aes.go
type ChachaCrypter struct {
	current  cipher.AEAD
	previous cipher.AEAD
}

func NewEncryptionProvider(key []byte, previousKey []byte) (sqlcrypter.Crypterer, error) {
	if len(key) != chacha20.KeySize {
		return nil, fmt.Errorf("invalid key length (want %v, got %v)", chacha20.KeySize, len(key))
	}

	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize an AEAD: %w", err)
	}

	var previous cipher.AEAD
	if previousKey != nil {
		if len(previousKey) != chacha20.KeySize {
			return nil, fmt.Errorf("invalid previous key length (want %v, got %v)", chacha20.KeySize, len(key))
		}

		previous, err = chacha20poly1305.New(previousKey)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize an AEAD: %w", err)
		}
	}

	return &ChachaCrypter{
		current: aead, previous: previous,
	}, nil
}

func (c *ChachaCrypter) Encrypt(w io.Writer, r io.Reader) error {
	src := new(bytes.Buffer)
	_, err := src.ReadFrom(r)
	if err != nil {
		return fmt.Errorf("failed to read from io.Reader: %w", err)
	}

	nonce, err := sqlcrypter.GenerateBytes(c.current.NonceSize())
	if err != nil {
		return fmt.Errorf("failed to generate 12-byte random nonce: %w", err)
	}

	ciphertext := c.current.Seal(nil, nonce, src.Bytes(), nil)

	// First 12 bytes will be the nonce, followed by the ciphertext
	w.Write(nonce)
	w.Write(ciphertext)

	return nil
}

func (a *ChachaCrypter) Decrypt(w io.Writer, r io.Reader) error {
	src := new(bytes.Buffer)
	n, err := src.ReadFrom(r)
	if err != nil {
		return fmt.Errorf("failed to read from io.Reader: %w", err)
	}

	// First 12 bytes is the nonce, followed by the ciphertext
	nonce := src.Next(chacha20poly1305.NonceSizeX)
	ciphertext := src.Next(int(n))

	// First attempt to decrypt using previous DEK if specified
	var attempted bool
	if a.previous != nil {
		if plaintext, err := a.previous.Open(nil, nonce, ciphertext, nil); err == nil {
			w.Write(plaintext)
			return nil
		}

		attempted = true
	}

	// Decrypt using the current DEK
	plaintext, err := a.current.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		if attempted {
			return fmt.Errorf("failed to decrypt ciphertext using current and previous DEK: %w", err)
		}

		return fmt.Errorf("failed to decrypt ciphertext using current DEK: %w", err)
	}

	w.Write(plaintext)

	return nil
}

var _ sqlcrypter.Crypterer = (*ChachaCrypter)(nil)
