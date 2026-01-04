package types

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
)

type CreateEventRegistrationParams struct {
	EventID    string `schema:"event_id"`
	GivenName  string `schema:"given_name"`
	FamilyName string `schema:"family_name"`
	Email      string `schema:"email"`
	Country    string `schema:"country"`
	Profession string `schema:"profession"`
	Company    string `schema:"company"`
}

var EmailValidationRegexp = regexp.MustCompile(`^[\w\-\.]+@([\w-]+\.)+[\w-]{2,}$`)

func (p *CreateEventRegistrationParams) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.Email, validation.Required, validation.Match(EmailValidationRegexp)),
		validation.Field(&p.Country, validation.Required, validation.Length(2, 2)),
	)
}
