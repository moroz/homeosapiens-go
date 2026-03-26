package types

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
)

var EmailValidationRegexp = regexp.MustCompile(`^[\w\-\.]+@([\w-]+\.)+[\w-]{2,}$`)

type OrderParams struct {
	Email               string `form:"email"`
	BillingGivenName    string `form:"billing_given_name"`
	BillingFamilyName   string `form:"billing_family_name"`
	BillingPhone        string `form:"billing_phone"`
	BillingAddressLine1 string `form:"billing_address_line1"`
	BillingAddressLine2 string `form:"billing_address_line2"`
	BillingCity         string `form:"billing_city"`
	BillingPostalCode   string `form:"billing_postal_code"`
	BillingCountry      string `form:"billing_country"`
}

func (p *OrderParams) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.Email, validation.Required, validation.Match(EmailValidationRegexp)),
	)
}
