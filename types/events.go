package types

import (
	"regexp"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/shopspring/decimal"
)

type UpdateEventInput struct {
	EventType string

	TitleEn    string
	TitlePl    string
	SubtitleEn *string
	SubtitlePl *string
	Slug       string

	DescriptionEn string
	DescriptionPl string

	// Pricing
	Price    *string
	Currency *string

	StartsAt time.Time
	EndsAt   time.Time

	// HostIds IDs of hosts associated with the event.
	HostIds []uuid.UUID
}

// PatchEventInput describes a selective update of an event. Every field is
// optional: absent fields are left untouched, and fields backed by a nullable
// column may be cleared by passing an explicit null.
type PatchEventInput struct {
	TitleEn Optional[string] `json:"titleEn"`
	TitlePl Optional[string] `json:"titlePl"`

	SubtitleEn Optional[string] `json:"subtitleEn"`
	SubtitlePl Optional[string] `json:"subtitlePl"`

	DescriptionEn Optional[string] `json:"descriptionEn"`
	DescriptionPl Optional[string] `json:"descriptionPl"`
}

func (p *PatchEventInput) Validate() error {
	return validation.ValidateStruct(p,
		// title_en and title_pl are NOT NULL: they may be changed, never cleared.
		validation.Field(&p.TitleEn, NotBlankWhenSet),
		validation.Field(&p.TitlePl, NotBlankWhenSet),
	)
}

type CreateEventInput struct {
	EventType string

	TitleEn    string
	TitlePl    string
	SubtitleEn *string
	SubtitlePl *string
	Slug       string

	DescriptionEn *string
	DescriptionPl *string

	// Pricing
	Price    *decimal.Decimal
	Currency *string

	// HostIds IDs of hosts associated with the event.
	HostIds []uuid.UUID

	StartsAt time.Time
	EndsAt   time.Time

	IsVirtual        bool
	VenueNameEn      *string
	VenueNamePl      *string
	VenueStreet      *string
	VenueCityEn      *string
	VenueCityPl      *string
	VenuePostalCode  *string
	VenueCountryCode *string
}

func (p *CreateEventInput) Validate() error {
	hasPrice := p.Price != nil && !p.Price.Equal(decimal.Zero)

	return validation.ValidateStruct(p,
		validation.Field(&p.EventType, validation.Required, validation.In("seminar", "webinar")),

		validation.Field(&p.TitlePl, validation.Required),
		validation.Field(&p.TitleEn, validation.Required),
		validation.Field(&p.Slug, validation.Required, validation.Match(slugRegexp)),

		validation.Field(&p.DescriptionEn, validation.Required),
		validation.Field(&p.DescriptionPl, validation.Required),

		// Currency is required (and constrained) only when the event carries a price.
		validation.Field(&p.Currency,
			validation.When(hasPrice,
				validation.Required,
				validation.In("PLN", "EUR"))),

		validation.Field(&p.StartsAt, validation.Required),
		validation.Field(&p.EndsAt, validation.Required, validation.Min(p.StartsAt).Exclusive().
			Error("must be after the start time")),

		// Physical events must carry a venue; virtual events need none.
		validation.Field(&p.VenueNameEn, validation.Required.When(!p.IsVirtual)),
		validation.Field(&p.VenueNamePl, validation.Required.When(!p.IsVirtual)),
		validation.Field(&p.VenueStreet, validation.Required.When(!p.IsVirtual)),
		validation.Field(&p.VenueCityEn, validation.Required.When(!p.IsVirtual)),
		validation.Field(&p.VenueCityPl, validation.Required.When(!p.IsVirtual)),
		validation.Field(&p.VenuePostalCode, validation.Required.When(!p.IsVirtual)),
		validation.Field(&p.VenueCountryCode, validation.Required.When(!p.IsVirtual),
			validation.Length(2, 2)),
	)
}

var slugRegexp = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

type EventDetailsDto struct {
	*queries.Event
	Product           *queries.Product
	Prices            []*queries.ProductPrice
	Hosts             []*queries.ListHostsForEventsRow
	EventRegistration *queries.EventRegistration
	RegistrationCount int
	CountInCart       int
}

func (d *EventDetailsDto) IsFree() bool {
	return d.ProductID == nil || d.Product.BasePriceAmount.Equal(decimal.Zero)
}
