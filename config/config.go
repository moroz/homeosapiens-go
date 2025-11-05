package config

import (
	"encoding/base64"
	"log"
	"os"
)

func MustGetenv(name string) string {
	val := os.Getenv(name)
	if val == "" {
		log.Fatalf(`FATAL: Environment variable %s is not set!`, name)
	}
	return val
}

func MustGetenvBase64(name string) []byte {
	val := MustGetenv(name)
	bytes, err := base64.StdEncoding.DecodeString(val)
	if err != nil {
		log.Fatalf(`FATAL: Failed to decode environment variable %s from Base64`)
	}
	return bytes
}

var DatabaseUrl = MustGetenv("DATABASE_URL")
var SecretKeyBase = MustGetenvBase64("SECRET_KEY_BASE")
