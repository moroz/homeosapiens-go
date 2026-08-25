package types

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// HostInput carries the editable fields of a host. Create and update take the
// same shape: a host is a small, flat record with no derived state.
type HostInput struct {
	Salutation *string `json:"salutation"`
	GivenName  string  `json:"givenName"`
	FamilyName string  `json:"familyName"`
	Country    *string `json:"country"`
}

func (p *HostInput) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.GivenName, validation.Required),
		validation.Field(&p.FamilyName, validation.Required),
		validation.Field(&p.Country, is.CountryCode2),
	)
}
