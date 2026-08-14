package auth

import (
	"errors"

	"github.com/byteness/keyring"
)

type KeyringClient struct {
	ring keyring.Keyring
}

func NewClient(keychain, service string) (*KeyringClient, error) {
	ring, err := keyring.Open(keyring.Config{
		KeychainName: keychain,
		ServiceName:  service,
	})
	if err != nil {
		return nil, err
	}
	return &KeyringClient{ring}, nil
}

func (c *KeyringClient) FetchAPIToken(key string) (string, error) {
	item, err := c.ring.Get(key)
	if errors.Is(err, keyring.ErrKeyNotFound) {
		return "", nil
	}
	return string(item.Data), err
}

func (c *KeyringClient) SetAPIToken(key string, value string) error {
	return c.ring.Set(keyring.Item{
		Key:  key,
		Data: []byte(value),
	})
}
