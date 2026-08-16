package types

import validation "github.com/go-ozzo/ozzo-validation/v4"

type CreateBlogPostInput struct {
	Title    string
	Slug     string
	Language string
	Body     *string
}

func (p *CreateBlogPostInput) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.Title, validation.Required),
		validation.Field(&p.Title, validation.Required, validation.Match(slugRegexp)),
		validation.Field(&p.Language, validation.In("en", "pl")),
	)
}
