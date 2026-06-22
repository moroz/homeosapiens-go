package email_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	sqlcrypter "github.com/bincyber/go-sqlcrypter"
	"github.com/google/uuid"
	"github.com/moroz/homeosapiens-go/db/queries"
	projI18n "github.com/moroz/homeosapiens-go/i18n"
	emailtmpl "github.com/moroz/homeosapiens-go/tmpl/email"
	"github.com/moroz/homeosapiens-go/types"
	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustInitBundle(t *testing.T) *goi18n.Bundle {
	t.Helper()
	bundle, err := projI18n.InitBundle()
	require.NoError(t, err)
	return bundle
}

func testUserToken() *types.UserTokenDTO {
	return &types.UserTokenDTO{
		UserToken: &queries.UserToken{
			ID:         uuid.New(),
			UserID:     uuid.New(),
			Context:    "email_verification",
			Token:      []byte("hashed-token"),
			InsertedAt: time.Now(),
			ValidUntil: time.Now().Add(24 * time.Hour),
		},
		User: &queries.User{
			ID:         uuid.New(),
			Email:      sqlcrypter.NewEncryptedBytes("user@example.com"),
			GivenName:  sqlcrypter.NewEncryptedBytes("John"),
			FamilyName: sqlcrypter.NewEncryptedBytes("Doe"),
		},
		PlaintextToken: []byte("raw-plaintext-token"),
	}
}

func testOrder() *types.OrderDTO {
	postal := sqlcrypter.NewEncryptedBytes("00-001")
	return &types.OrderDTO{
		Order: &queries.Order{
			ID:                  uuid.New(),
			UserID:              uuid.New(),
			OrderNumber:         42,
			GrandTotal:          decimal.NewFromFloat(99.99),
			Currency:            "PLN",
			BillingGivenName:    sqlcrypter.NewEncryptedBytes("Jane"),
			BillingFamilyName:   sqlcrypter.NewEncryptedBytes("Smith"),
			Email:               sqlcrypter.NewEncryptedBytes("jane@example.com"),
			BillingCity:         sqlcrypter.NewEncryptedBytes("Warsaw"),
			BillingAddressLine1: sqlcrypter.NewEncryptedBytes("ul. Marszałkowska 1"),
			BillingPostalCode:   &postal,
			BillingCountry:      "PL",
			InsertedAt:          time.Now(),
			UpdatedAt:           time.Now(),
		},
		LineItems: []*queries.OrderLineItem{
			{
				ID:                   uuid.New(),
				OrderID:              uuid.New(),
				ProductTitle:         "Advanced Homeopathy Seminar",
				ProductPriceAmount:   decimal.NewFromFloat(49.99),
				ProductPriceCurrency: "PLN",
				Quantity:             1,
				InsertedAt:           time.Now(),
				UpdatedAt:            time.Now(),
			},
		},
	}
}

func TestUserEmailVerificationTemplate(t *testing.T) {
	bundle := mustInitBundle(t)

	locales := []struct {
		lang    string
		title   string
		ctaText string
	}{
		{"en", "Please confirm your email address", "Verify email address"},
		{"pl", "Potwierdź swój adres e-mail", "Potwierdzam rejestrację"},
	}

	for _, lc := range locales {
		t.Run(lc.lang, func(t *testing.T) {
			token := testUserToken()
			props := &emailtmpl.UserEmailVerificationEmailProps{
				LayoutProps: &emailtmpl.LayoutProps{
					Title:     lc.title,
					Language:  lc.lang,
					Localizer: goi18n.NewLocalizer(bundle, lc.lang),
				},
				UserToken: token,
			}

			var buf bytes.Buffer
			err := emailtmpl.UserEmailVerificationTemplate.Execute(&buf, props)
			require.NoError(t, err)

			html := buf.String()
			assert.True(t, strings.HasPrefix(html, "<!DOCTYPE html>"))
			assert.Contains(t, html, lc.title)
			assert.Contains(t, html, token.VerifyEmailURL())
			assert.Contains(t, html, lc.ctaText)
		})
	}
}

func TestUserPasswordResetTemplate(t *testing.T) {
	bundle := mustInitBundle(t)

	locales := []struct {
		lang    string
		title   string
		ctaText string
	}{
		{"en", "Password recovery", "Change my password"},
		{"pl", "Odzyskiwanie hasła", "Zmień hasło"},
	}

	for _, lc := range locales {
		t.Run(lc.lang, func(t *testing.T) {
			token := testUserToken()
			props := &emailtmpl.UserPasswordResetEmailProps{
				LayoutProps: &emailtmpl.LayoutProps{
					Title:     lc.title,
					Language:  lc.lang,
					Localizer: goi18n.NewLocalizer(bundle, lc.lang),
				},
				UserToken: token,
			}

			var buf bytes.Buffer
			err := emailtmpl.UserPasswordResetTemplate.Execute(&buf, props)
			require.NoError(t, err)

			html := buf.String()
			assert.True(t, strings.HasPrefix(html, "<!DOCTYPE html>"))
			assert.Contains(t, html, lc.title)
			assert.Contains(t, html, token.ResetPasswordURL())
			assert.Contains(t, html, lc.ctaText)
		})
	}
}

func TestOrderConfirmationTemplate(t *testing.T) {
	bundle := mustInitBundle(t)

	locales := []struct {
		lang           string
		heading        string
		orderNumberStr string
	}{
		{"en", "Order Confirmation", "#42"},
		{"pl", "Potwierdzenie zamówienia", "nr&nbsp;42"},
	}

	for _, lc := range locales {
		t.Run(lc.lang, func(t *testing.T) {
			order := testOrder()
			props := &emailtmpl.OrderEmailProps{
				LayoutProps: &emailtmpl.LayoutProps{
					Title:     lc.heading,
					Language:  lc.lang,
					Localizer: goi18n.NewLocalizer(bundle, lc.lang),
				},
				Order: order,
			}

			var buf bytes.Buffer
			err := emailtmpl.OrderConfirmationTemplate.Execute(&buf, props)
			require.NoError(t, err)

			html := buf.String()
			assert.True(t, strings.HasPrefix(html, "<!DOCTYPE html>"))
			assert.Contains(t, html, lc.heading)
			assert.Contains(t, html, "Advanced Homeopathy Seminar")
			assert.Contains(t, html, lc.orderNumberStr)
		})
	}
}

func testEventRegistrationEmailDTO(lang string) *types.EventRegistrationEmailDTO {
	startsAt := time.Date(2026, 9, 15, 14, 0, 0, 0, time.UTC)
	endsAt := time.Date(2026, 9, 15, 17, 0, 0, 0, time.UTC)
	venueName := "Centrum Konferencyjne"
	venueStreet := "ul. Marszałkowska 1"
	venueCity := "Warszawa"
	postalCode := "00-001"
	return &types.EventRegistrationEmailDTO{
		Event: &queries.Event{
			ID:           uuid.New(),
			TitleEn:      "Advanced Homeopathy Seminar",
			TitlePl:      "Zaawansowane seminarium homeopatyczne",
			Slug:         "advanced-homeopathy-seminar",
			StartsAt:     startsAt,
			EndsAt:       endsAt,
			IsVirtual:    false,
			EventType:    queries.EventTypeSeminar,
			VenueNameEn:  &venueName,
			VenueNamePl:  &venueName,
			VenueStreet:  &venueStreet,
			VenueCityEn:  &venueCity,
			VenueCityPl:  &venueCity,
			VenuePostalCode: &postalCode,
			DescriptionEn: "A hands-on seminar.",
		},
		User: &queries.User{
			ID:              uuid.New(),
			GivenName:       sqlcrypter.NewEncryptedBytes("Anna"),
			FamilyName:      sqlcrypter.NewEncryptedBytes("Kowalska"),
			Email:           sqlcrypter.NewEncryptedBytes("anna@example.com"),
			PreferredLocale: queries.Locale(lang),
		},
	}
}

func TestEventRegistrationConfirmationTemplate(t *testing.T) {
	bundle := mustInitBundle(t)

	locales := []struct {
		lang         string
		heading      string
		title        string
		formattedDate string
	}{
		{"en", "Registration Confirmed", "Advanced Homeopathy Seminar", "September 15, 2026 at 14:00 UTC"},
		{"pl", "Rejestracja potwierdzona", "Zaawansowane seminarium homeopatyczne", "15.09.2026, 14:00 UTC"},
	}

	for _, lc := range locales {
		t.Run(lc.lang, func(t *testing.T) {
			data := testEventRegistrationEmailDTO(lc.lang)
			props := &emailtmpl.EventRegistrationEmailProps{
				LayoutProps: &emailtmpl.LayoutProps{
					Title:     lc.heading,
					Language:  lc.lang,
					Localizer: goi18n.NewLocalizer(bundle, lc.lang),
				},
				Data: data,
			}

			var buf bytes.Buffer
			err := emailtmpl.EventRegistrationConfirmationTemplate.Execute(&buf, props)
			require.NoError(t, err)

			html := buf.String()
			assert.True(t, strings.HasPrefix(html, "<!DOCTYPE html>"))
			assert.Contains(t, html, lc.heading)
			assert.Contains(t, html, lc.title)
			assert.Contains(t, html, lc.formattedDate)
			assert.Contains(t, html, data.EventURL())
		})
	}
}

func TestEventRegistrationConfirmationTemplate_VirtualEvent(t *testing.T) {
	bundle := mustInitBundle(t)
	data := testEventRegistrationEmailDTO("en")
	data.Event.IsVirtual = true

	props := &emailtmpl.EventRegistrationEmailProps{
		LayoutProps: &emailtmpl.LayoutProps{
			Title:     "Registration Confirmed",
			Language:  "en",
			Localizer: goi18n.NewLocalizer(bundle, "en"),
		},
		Data: data,
	}

	var buf bytes.Buffer
	err := emailtmpl.EventRegistrationConfirmationTemplate.Execute(&buf, props)
	require.NoError(t, err)

	html := buf.String()
	assert.Contains(t, html, "Online event")
}

func TestPaymentConfirmationTemplate(t *testing.T) {
	bundle := mustInitBundle(t)

	locales := []struct {
		lang    string
		heading string
	}{
		{"en", "Payment Confirmed"},
		{"pl", "Potwierdzenie płatności"},
	}

	for _, lc := range locales {
		t.Run(lc.lang, func(t *testing.T) {
			order := testOrder()
			props := &emailtmpl.OrderEmailProps{
				LayoutProps: &emailtmpl.LayoutProps{
					Title:     lc.heading,
					Language:  lc.lang,
					Localizer: goi18n.NewLocalizer(bundle, lc.lang),
				},
				Order: order,
			}

			var buf bytes.Buffer
			err := emailtmpl.PaymentConfirmationTemplate.Execute(&buf, props)
			require.NoError(t, err)

			html := buf.String()
			assert.True(t, strings.HasPrefix(html, "<!DOCTYPE html>"))
			assert.Contains(t, html, lc.heading)
			assert.Contains(t, html, "Advanced Homeopathy Seminar")
		})
	}
}
