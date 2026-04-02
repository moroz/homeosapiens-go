package handlers_test

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bincyber/go-sqlcrypter"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/internal/crypto"
	"github.com/moroz/homeosapiens-go/internal/mailer"
	"github.com/moroz/homeosapiens-go/services/mocks"
	"github.com/moroz/homeosapiens-go/web/router"
	"github.com/moroz/homeosapiens-go/web/session"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func uniqueEmail() string {
	unique := make([]byte, 4)
	_, _ = rand.Read(unique)
	return "user-" + hex.EncodeToString(unique) + "@example.com"
}

func TestStripeWebhook(t *testing.T) {
	db, err := initDB(t.Context())
	require.NoError(t, err)
	defer db.Close()

	store, err := session.NewStore(config.SessionKey)
	require.NoError(t, err)

	mail := mailer.MockMailer()

	stripe := mocks.NewMockStripeService(t)
	r := router.Router(db, store, stripe, mail)
	_ = r

	email := uniqueEmail()

	user, err := queries.New(db).InsertUser(t.Context(), &queries.InsertUserParams{
		Email:      sqlcrypter.NewEncryptedBytes(email),
		EmailHash:  crypto.HashEmail(email),
		GivenName:  sqlcrypter.NewEncryptedBytes("John"),
		FamilyName: sqlcrypter.NewEncryptedBytes("Smith"),
		Country:    new("IE"),
	})
	require.NoError(t, err)

	order, err := queries.New(db).InsertOrder(t.Context(), &queries.InsertOrderParams{
		UserID:              user.ID,
		GrandTotal:          decimal.NewFromFloat(42.),
		Currency:            "PLN",
		BillingGivenName:    sqlcrypter.NewEncryptedBytes("John"),
		BillingFamilyName:   sqlcrypter.NewEncryptedBytes("Smith"),
		BillingCity:         sqlcrypter.NewEncryptedBytes("Some City"),
		BillingPostalCode:   new(sqlcrypter.NewEncryptedBytes("12345")),
		BillingCountry:      "IE",
		Email:               sqlcrypter.NewEncryptedBytes(email),
		BillingAddressLine1: sqlcrypter.NewEncryptedBytes("42, Some Street"),
	})

	require.NoError(t, err)
	require.NotNil(t, order)

	cs := mocks.GenerateCheckoutSession()

	order, err = queries.New(db).StoreCheckoutSessionIDOnOrder(t.Context(), &queries.StoreCheckoutSessionIDOnOrderParams{
		StripeCheckoutSessionID: &cs.ID,
		ID:                      order.ID,
	})
	require.NoError(t, err)

	assert.Nil(t, order.PaidAt)

	stripe.EXPECT().DecodeWebhook(mock.Anything, mock.Anything).Return(cs, nil)

	req, _ := http.NewRequest("POST", "/webhooks/stripe", bytes.NewBufferString("{}"))
	req.Header.Add("Stripe-Signature", "PHONY")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	resp := rr.Result()
	assert.Equal(t, 200, resp.StatusCode)
}
