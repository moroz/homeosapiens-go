package config

import (
	"crypto/sha512"
	"encoding/base64"
	"io"
	"log"
	"os"

	"golang.org/x/crypto/hkdf"
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
		log.Fatalf(`FATAL: Failed to decode environment variable %s from Base64`, name)
	}
	return bytes
}

func GetEnvWithDefault(name, defaultValue string) string {
	val := os.Getenv(name)
	if val == "" {
		return defaultValue
	}
	return val
}

func MustDeriveKey(base []byte, info string, lengthInBytes int) []byte {
	kdf := hkdf.New(sha512.New, base, nil, []byte(info))
	buf := make([]byte, lengthInBytes)
	if _, err := io.ReadFull(kdf, buf); err != nil {
		log.Fatalf("Failed to derive key (info: %s): %s", info, err)
	}
	return buf
}

func RequireInProduction(name string) string {
	if IsProd {
		return MustGetenv(name)
	}
	return os.Getenv(name)
}

var AppPort = GetEnvWithDefault("PORT", "3000")
var DatabaseUrl = MustGetenv("DATABASE_URL")
var SecretKeyBase = MustGetenvBase64("SECRET_KEY_BASE")
var SessionKey = MustDeriveKey(SecretKeyBase, "Sessions", 32)
var IsProd = os.Getenv("GO_ENV") == "prod"
var GoogleClientId = RequireInProduction("GOOGLE_CLIENT_ID")
var GoogleClientSecret = RequireInProduction("GOOGLE_CLIENT_SECRET")

const AssetCdnBaseUrl = "https://d3n1g0yg3ja4p3.cloudfront.net"
const SessionCookieName = "_hs_session"
const SessionContextName = "session"
const CurrentUserContextName = "current_user"
