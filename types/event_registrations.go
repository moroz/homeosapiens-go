package types

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
)

type CreateEventRegistrationParams struct {
	EventID    string `form:"event_id"`
	GivenName  string `form:"given_name"`
	FamilyName string `form:"family_name"`
	Email      string `form:"email"`
	Country    string `form:"country"`
	Profession string `form:"profession"`
	Company    string `form:"company"`
}

var EmailValidationRegexp = regexp.MustCompile(`^[\w\-\.]+@([\w-]+\.)+[\w-]{2,}$`)

func (p *CreateEventRegistrationParams) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.Email, validation.Required, validation.Match(EmailValidationRegexp)),
		validation.Field(&p.Country, validation.Required, validation.Length(2, 2)),
	)
}
