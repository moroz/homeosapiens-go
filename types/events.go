package types

import (
	"regexp"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/shopspring/decimal"
)

// PatchEventInput describes a selective update of an event. Every field is
// optional: absent fields are left untouched, and fields backed by a nullable
// column may be cleared by passing an explicit null.
type PatchEventInput struct {
	TitleEn Optional[string] `json:"titleEn"`
	TitlePl Optional[string] `json:"titlePl"`
	Slug    Optional[string] `json:"slug"`

	SubtitleEn Optional[string] `json:"subtitleEn"`
	SubtitlePl Optional[string] `json:"subtitlePl"`

	DescriptionEn Optional[string] `json:"descriptionEn"`
	DescriptionPl Optional[string] `json:"descriptionPl"`

	StartsAt Optional[time.Time] `json:"startsAt"`
	EndsAt   Optional[time.Time] `json:"endsAt"`

	// Pricing lives on the event's product, not on the event itself. A null or
	// zero price makes the event free without discarding the product, so that
	// order and cart line items referencing it keep resolving.
	Price    Optional[decimal.Decimal] `json:"price"`
	Currency Optional[string]          `json:"currency"`

	// HostIds replaces the event's hosts wholesale, in the given order. Hosts
	// are positional, so there is no meaningful way to patch one in isolation.
	HostIds Optional[[]uuid.UUID] `json:"hostIds"`
}

// IsEmpty reports whether the payload carries no changes at all, in which case
// the record should be left alone rather than having its updated_at bumped.
func (p *PatchEventInput) IsEmpty() bool {
	return !(p.TitleEn.Set || p.TitlePl.Set || p.Slug.Set ||
		p.SubtitleEn.Set || p.SubtitlePl.Set ||
		p.DescriptionEn.Set || p.DescriptionPl.Set ||
		p.Price.Set || p.Currency.Set || p.HostIds.Set)
}

func (p *PatchEventInput) Validate() error {
	return validation.ValidateStruct(p,
		// title_en, title_pl and slug are NOT NULL: they may be changed, never cleared.
		validation.Field(&p.TitleEn, NotBlankWhenSet),
		validation.Field(&p.TitlePl, NotBlankWhenSet),

		// Slug uniqueness cannot be checked here without racing another writer;
		// the service catches the constraint violation instead.
		validation.Field(&p.Slug, NotBlankWhenSet, WhenSet[string](validation.Match(slugRegexp))),

		validation.Field(&p.Price, WhenSet[decimal.Decimal](validation.By(nonNegative))),
		validation.Field(&p.Currency, NotBlankWhenSet, WhenSet[string](validation.In("PLN", "EUR"))),

		validation.Field(&p.HostIds, WhenSet[[]uuid.UUID](validation.By(distinctIDs))),

		validation.Field(&p.StartsAt, NotBlankWhenSet),
		validation.Field(&p.EndsAt, NotBlankWhenSet, WhenSet[time.Time](validation.Min(p.StartsAt.Value).Exclusive().Error("must be after start time"))),
	)
}

func nonNegative(value any) error {
	amount, ok := value.(decimal.Decimal)
	if !ok || !amount.IsNegative() {
		return nil
	}

	return validation.NewError("min", "must not be negative")
}

func distinctIDs(value any) error {
	ids, ok := value.([]uuid.UUID)
	if !ok {
		return nil
	}

	seen := make(map[uuid.UUID]struct{}, len(ids))
	for _, id := range ids {
		if _, duplicate := seen[id]; duplicate {
			return validation.NewError("distinct", "must not contain duplicates")
		}
		seen[id] = struct{}{}
	}

	return nil
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
