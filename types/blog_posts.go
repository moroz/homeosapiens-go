package types

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type CreateBlogPostInput struct {
	Title    string  `json:"title"`
	Slug     string  `json:"slug"`
	Language string  `json:"language"`
	Body     *string `json:"body"`
}

func (p *CreateBlogPostInput) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.Title, validation.Required),
		validation.Field(&p.Slug, validation.Required, validation.Match(slugRegexp)),
		validation.Field(&p.Language, validation.In("en", "pl")),
	)
}

type UpdateBlogPostInput struct {
	Title    string  `json:"title"`
	Slug     string  `json:"slug"`
	Language string  `json:"language"`
	Body     *string `json:"body"`
}

func (p *UpdateBlogPostInput) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.Title, validation.Required),
		validation.Field(&p.Slug, validation.Required, validation.Match(slugRegexp)),
		validation.Field(&p.Language, validation.In("en", "pl")),
	)
}

type PublishBlogPostValidation struct {
	Body        *string    `json:"body"`
	PublishedAt *time.Time `json:"publishedAt"`
}

func (p *PublishBlogPostValidation) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.Body, validation.Required),
		validation.Field(&p.PublishedAt, validation.Nil.Error("post is already published")),
	)
}

type UnpublishBlogPostValidation struct {
	PublishedAt *time.Time `json:"publishedAt"`
}

func (p *UnpublishBlogPostValidation) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.PublishedAt, validation.Required.Error("post is not published")),
	)
}
