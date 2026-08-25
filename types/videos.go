package types

import (
	"time"

	"github.com/google/uuid"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/shopspring/decimal"
)

type VideoGroupListDTO struct {
	*queries.VideoGroup
	HasAccess     bool
	MinRecordedOn *time.Time
	MaxRecordedOn *time.Time
	// Price and Currency are the price of the group's product, so that a locked
	// group can be sold on the spot. Both are nil for a free group.
	Price    *decimal.Decimal
	Currency *string
	// CountInCart is how many times the group's product sits in the visitor's
	// cart, mirroring how paid events render their add-to-cart button.
	CountInCart int
}

type VideoGroupDetailsDTO struct {
	*queries.VideoGroup
	HasAccess   bool
	Videos      []*queries.Video
	Price       *decimal.Decimal
	Currency    *string
	CountInCart int
}

type VideoDetailsDTO struct {
	*queries.Video
	HasAccess bool
	Sources   []*queries.VideoSource
}

type VideoListDTO struct {
	*queries.Video
	Hosts []*queries.Host
}

// CreateVideoGroupInput describes a new video group. The json tags are what
// ozzo reports validation errors under, so they must match the property names of
// the VideoGroupInput schema in web/api/openapi.yaml.
type CreateVideoGroupInput struct {
	TitleEn string `json:"titleEn"`
	TitlePl string `json:"titlePl"`
	Slug    string `json:"slug"`

	// Pricing lives on the group's product. A null or zero price leaves the group
	// free, i.e. without a product at all.
	Price    *decimal.Decimal `json:"price"`
	Currency *string          `json:"currency"`
}

func (p *CreateVideoGroupInput) Validate() error {
	hasPrice := p.Price != nil && !p.Price.Equal(decimal.Zero)

	return validation.ValidateStruct(p,
		validation.Field(&p.TitleEn, validation.Required),
		validation.Field(&p.TitlePl, validation.Required),
		validation.Field(&p.Slug, validation.Required, validation.Match(slugRegexp)),
		validation.Field(&p.Price, validation.By(nonNegativePtr)),
		validation.Field(&p.Currency,
			validation.When(hasPrice,
				validation.Required,
				validation.In("PLN", "EUR"))),
	)
}

// PatchVideoGroupInput describes a selective update of a video group. Absent
// fields are left untouched.
type PatchVideoGroupInput struct {
	TitleEn Optional[string] `json:"titleEn"`
	TitlePl Optional[string] `json:"titlePl"`
	Slug    Optional[string] `json:"slug"`

	Price    Optional[decimal.Decimal] `json:"price"`
	Currency Optional[string]          `json:"currency"`
}

// IsEmpty reports whether the payload carries no changes at all, in which case
// the record should be left alone rather than having its updated_at bumped.
func (p *PatchVideoGroupInput) IsEmpty() bool {
	return !(p.TitleEn.Set || p.TitlePl.Set || p.Slug.Set || p.Price.Set || p.Currency.Set)
}

func (p *PatchVideoGroupInput) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.TitleEn, NotBlankWhenSet),
		validation.Field(&p.TitlePl, NotBlankWhenSet),

		// Slug uniqueness cannot be checked here without racing another writer;
		// the service catches the constraint violation instead.
		validation.Field(&p.Slug, NotBlankWhenSet, WhenSet[string](validation.Match(slugRegexp))),

		validation.Field(&p.Price, WhenSet[decimal.Decimal](validation.By(nonNegative))),
		validation.Field(&p.Currency, NotBlankWhenSet, WhenSet[string](validation.In("PLN", "EUR"))),
	)
}

func nonNegativePtr(value any) error {
	amount, ok := value.(*decimal.Decimal)
	if !ok || amount == nil {
		return nil
	}

	return nonNegative(*amount)
}

func (v *VideoGroupListDTO) IsPremium() bool {
	return v.VideoGroup.ProductID != nil
}

func (v *VideoGroupDetailsDTO) IsPremium() bool {
	return v.VideoGroup.ProductID != nil
}

// UpdateVideoInput carries the editable fields of a video. Provider, youtube id,
// duration and thumbnails are set by the import script and stay read-only here,
// so an update replaces everything an admin is allowed to touch.
type UpdateVideoInput struct {
	TitleEn       string     `json:"titleEn"`
	TitlePl       string     `json:"titlePl"`
	Slug          string     `json:"slug"`
	DescriptionEn *string    `json:"descriptionEn"`
	DescriptionPl *string    `json:"descriptionPl"`
	RecordedOn    *time.Time `json:"recordedOn"`
	IsPublic      bool       `json:"isPublic"`
	HostID        *uuid.UUID `json:"hostId"`
}

func (p *UpdateVideoInput) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.TitleEn, validation.Required),
		validation.Field(&p.TitlePl, validation.Required),

		// Slug uniqueness cannot be checked here without racing another writer;
		// the service catches the constraint violation instead.
		validation.Field(&p.Slug, validation.Required, validation.Match(slugRegexp)),
	)
}
