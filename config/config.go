package config

import (
	"crypto/sha512"
	"encoding/base64"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/alexedwards/argon2id"
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

func MustParsePort(val string) int {
	parsed, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		log.Fatalf("Failed to parse value %v as port number: %s", val, err)
	}
	return int(parsed)
}

func RequireInProduction(name string, defaultValue string) string {
	if IsProd {
		return MustGetenv(name)
	}
	return GetEnvWithDefault(name, defaultValue)
}

func ResolveArgon2Params() *argon2id.Params {
	if IsProd {
		return &argon2id.Params{
			Memory:      64 * 1024,
			Iterations:  1,
			Parallelism: 1,
			SaltLength:  16,
			KeyLength:   16,
		}
	}

	return &argon2id.Params{
		Memory:      16 * 1024,
		Iterations:  1,
		Parallelism: 1,
		SaltLength:  16,
		KeyLength:   16,
	}
}

var AppPort = GetEnvWithDefault("PORT", "3000")
var DatabaseUrl = MustGetenv("DATABASE_URL")
var SecretKeyBase = MustGetenvBase64("SECRET_KEY_BASE")
var SessionKey = MustDeriveKey(SecretKeyBase, "Sessions", 32)
var IsProd = os.Getenv("GO_ENV") == "prod"
var IsTest = os.Getenv("GO_ENV") == "test"
var GoogleClientId = RequireInProduction("GOOGLE_CLIENT_ID", "")
var GoogleClientSecret = RequireInProduction("GOOGLE_CLIENT_SECRET", "")
var PublicUrl = RequireInProduction("PUBLIC_URL", "http://localhost:3000")
var StripeSecretKey = RequireInProduction("STRIPE_SECRET_KEY", "")
var StripePublishableKey = RequireInProduction("STRIPE_PUBLISHABLE_KEY", "")
var StripeWebhookSigningSecret = RequireInProduction("STRIPE_WEBHOOK_SIGNING_SECRET", "")

var SMTPHost = RequireInProduction("SMTP_HOST", "localhost")
var SMTPPort = MustParsePort(RequireInProduction("SMTP_PORT", "1025"))
var SMTPUsername = RequireInProduction("SMTP_USERNAME", "smtpuser")
var SMTPPassword = RequireInProduction("SMTP_PASSWORD", "smtppassword")

const AssetCdnBaseUrl = "https://d3n1g0yg3ja4p3.cloudfront.net"
const SessionCookieName = "_hs_session"

const OAuth2SessionKey = "auth_state"
const RedirectBackUrlSessionKey = "redirect_back"
const FlashSessionKey = "_flash"
const CartIdSessionKey = "cart_id"

const MinPasswordLength = 8
const MaxPasswordLength = 128

var DatabaseEncryptionKey = MustDeriveKey(SecretKeyBase, "ColumnLevelEncryption", 32)
var DatabaseHMACKey = MustDeriveKey(SecretKeyBase, "DatabaseHMAC", 32)
