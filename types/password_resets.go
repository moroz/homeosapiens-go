package types

import validation "github.com/go-ozzo/ozzo-validation/v4"

type UpdateUserPasswordParams struct {
	Token                string `form:"-"`
	Password             string `form:"password"`
	PasswordConfirmation string `form:"password_confirmation"`
}

func (p *UpdateUserPasswordParams) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.Password, validation.Required, validation.Length(8, 128)),
		validation.Field(&p.PasswordConfirmation, validation.In(p.Password).ErrorObject(
			validation.NewError("validation_password_confirmation", "does not match"),
		)),
	)
}
