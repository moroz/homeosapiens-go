package types

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
)

type CreateEventRegistrationParams struct {
	GivenName  string
	FamilyName string
	Email      string
	Country    string
	Profession string
}

var EmailValidationRegexp = regexp.MustCompile(`^[\w\-\.]+@([\w-]+\.)+[\w-]{2,}$`)

func (p *CreateEventRegistrationParams) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.Email, validation.Required, validation.Match(EmailValidationRegexp)),
		validation.Field(&p.Country, validation.Required, validation.Length(2, 2)),
	)
}
