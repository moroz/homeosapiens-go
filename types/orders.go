package types

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/stripe/stripe-go/v84"
)

var EmailValidationRegexp = regexp.MustCompile(`^[\w\-\.]+@([\w-]+\.)+[\w-]{2,}$`)

type OrderParams struct {
	Email               string `form:"email" json:"email"`
	BillingGivenName    string `form:"billing_given_name" json:"billing_given_name""`
	BillingFamilyName   string `form:"billing_family_name" json:"billing_family_name"`
	BillingPhone        string `form:"billing_phone" json:"billing_phone"`
	BillingAddressLine1 string `form:"billing_address_line1" json:"billing_address_line1"`
	BillingAddressLine2 string `form:"billing_address_line2" json:"billing_address_line2"`
	BillingCity         string `form:"billing_city" json:"billing_city"`
	BillingPostalCode   string `form:"billing_postal_code" json:"billing_postal_code"`
	BillingCountry      string `form:"billing_country" json:"billing_country"`
}

func (p *OrderParams) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.Email, validation.Required, validation.Match(EmailValidationRegexp)),
		validation.Field(&p.BillingGivenName, validation.Required),
		validation.Field(&p.BillingFamilyName, validation.Required),
		validation.Field(&p.BillingAddressLine1, validation.Required),
		validation.Field(&p.BillingCity, validation.Required),
		validation.Field(&p.BillingCountry, validation.Required),
		validation.Field(&p.BillingPostalCode, validation.Required),
	)
}

type OrderDTO struct {
	*queries.Order
	LineItems       []*queries.OrderLineItem
	CheckoutSession *stripe.CheckoutSession
}
