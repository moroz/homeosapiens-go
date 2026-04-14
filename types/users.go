package types

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/moroz/homeosapiens-go/db/queries"
)

// SeedUserParams represents the trusted data of a user, coming from a trusted source, such as database data exports. It is intended for backfilling user data, or for seeding the database.
type SeedUserParams struct {
	GivenName      string
	FamilyName     string
	Email          string
	LicenceNumber  string
	Role           queries.UserRole
	Country        string
	Password       string
	EmailConfirmed bool
}

type RegisterUserParams struct {
	PreferredLocale      string `form:"locale" json:"locale"`
	GivenName            string `form:"given_name" json:"given_name"`
	FamilyName           string `form:"family_name" json:"family_name"`
	Email                string `form:"email" json:"email"`
	Password             string `form:"password" json:"password"`
	PasswordConfirmation string `form:"password_confirmation" json:"password_confirmation"`
}

func (p *RegisterUserParams) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.PreferredLocale, validation.Required, validation.In("pl", "en")),
		validation.Field(&p.GivenName, validation.Required),
		validation.Field(&p.FamilyName, validation.Required),
		validation.Field(&p.Email, validation.Required, is.Email),
		validation.Field(&p.Password, validation.Required, validation.Length(8, 128)),
		validation.Field(&p.PasswordConfirmation, validation.In(p.Password).ErrorObject(
			validation.NewError("validation_password_confirmation", "does not match"),
		)),
	)
}

type UserRegistrationDTO struct {
	*queries.User
	*UserTokenDTO
}
