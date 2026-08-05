package types

import (
	"regexp"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
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

	// IsVirtual and the venue fields are coupled: a physical event must carry a
	// full venue, which is enforced against the state the event will have once
	// the patch is applied (see ValidateVenue).
	IsVirtual Optional[bool] `json:"isVirtual"`

	// MeetingUrl is the join link (Zoom or similar) sent out with registration
	// confirmations and reminders. It is optional even for virtual events, which
	// may be scheduled before the meeting has been created.
	MeetingUrl Optional[string] `json:"meetingUrl"`

	VenueNameEn      Optional[string] `json:"venueNameEn"`
	VenueNamePl      Optional[string] `json:"venueNamePl"`
	VenueStreet      Optional[string] `json:"venueStreet"`
	VenueCityEn      Optional[string] `json:"venueCityEn"`
	VenueCityPl      Optional[string] `json:"venueCityPl"`
	VenuePostalCode  Optional[string] `json:"venuePostalCode"`
	VenueCountryCode Optional[string] `json:"venueCountryCode"`
}

// IsEmpty reports whether the payload carries no changes at all, in which case
// the record should be left alone rather than having its updated_at bumped.
func (p *PatchEventInput) IsEmpty() bool {
	return !(p.TitleEn.Set || p.TitlePl.Set || p.Slug.Set ||
		p.SubtitleEn.Set || p.SubtitlePl.Set ||
		p.DescriptionEn.Set || p.DescriptionPl.Set ||
		p.StartsAt.Set || p.EndsAt.Set ||
		p.Price.Set || p.Currency.Set || p.HostIds.Set ||
		p.IsVirtual.Set || p.MeetingUrl.Set || p.venueFieldsSet())
}

func (p *PatchEventInput) venueFieldsSet() bool {
	for _, field := range p.venueFields() {
		if field.value.Set {
			return true
		}
	}
	return false
}

// venueFields pairs every venue field with its payload property name and with
// the event field it is read back from, so that the emptiness check and the
// venue validation iterate the same list.
func (p *PatchEventInput) venueFields() []struct {
	name   string
	value  Optional[string]
	stored func(event *queries.Event) *string
} {
	return []struct {
		name   string
		value  Optional[string]
		stored func(event *queries.Event) *string
	}{
		{"venueNameEn", p.VenueNameEn, func(e *queries.Event) *string { return e.VenueNameEn }},
		{"venueNamePl", p.VenueNamePl, func(e *queries.Event) *string { return e.VenueNamePl }},
		{"venueStreet", p.VenueStreet, func(e *queries.Event) *string { return e.VenueStreet }},
		{"venueCityEn", p.VenueCityEn, func(e *queries.Event) *string { return e.VenueCityEn }},
		{"venueCityPl", p.VenueCityPl, func(e *queries.Event) *string { return e.VenueCityPl }},
		{"venuePostalCode", p.VenuePostalCode, func(e *queries.Event) *string { return e.VenuePostalCode }},
		{"venueCountryCode", p.VenueCountryCode, func(e *queries.Event) *string { return e.VenueCountryCode }},
	}
}

// ValidateVenue checks the venue against the state the event will have once the
// patch is applied: a physical event must carry a full venue, whether the
// address comes from the payload or from the row being patched. The database
// enforces a subset of this with a check constraint, but a constraint violation
// would surface as a 500 rather than as a per-field error.
func (p *PatchEventInput) ValidateVenue(event *queries.Event) error {
	isVirtual := event.IsVirtual
	if p.IsVirtual.Set {
		isVirtual = p.IsVirtual.Value
	}

	if isVirtual {
		return nil
	}

	errs := validation.Errors{}
	for _, field := range p.venueFields() {
		value := field.stored(event)
		if field.value.Set {
			value = field.value.Ptr()
		}

		if value == nil || strings.TrimSpace(*value) == "" {
			errs[field.name] = validation.ErrRequired
		}
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
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

		// is_virtual is NOT NULL, so it may be flipped but never cleared.
		validation.Field(&p.IsVirtual, NotNullWhenSet[bool]()),

		// meeting_url is nullable, so an explicit null clears the join link.
		validation.Field(&p.MeetingUrl, NotBlankWhenPresent, WhenSet[string](is.URL)),

		// The venue columns are nullable, so an explicit null is how a venue is
		// cleared; a present value must still be a real one. Whether the
		// resulting event is allowed to have no venue is decided by ValidateVenue.
		validation.Field(&p.VenueNameEn, NotBlankWhenPresent),
		validation.Field(&p.VenueNamePl, NotBlankWhenPresent),
		validation.Field(&p.VenueStreet, NotBlankWhenPresent),
		validation.Field(&p.VenueCityEn, NotBlankWhenPresent),
		validation.Field(&p.VenueCityPl, NotBlankWhenPresent),
		validation.Field(&p.VenuePostalCode, NotBlankWhenPresent),
		validation.Field(&p.VenueCountryCode, NotBlankWhenPresent,
			WhenSet[string](validation.Length(2, 2))),
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

// CreateEventInput describes a new event. The json tags are what ozzo reports
// validation errors under, so they must match the property names of the
// EventInput schema in web/api/openapi.yaml.
type CreateEventInput struct {
	EventType string `json:"eventType"`

	TitleEn    string  `json:"titleEn"`
	TitlePl    string  `json:"titlePl"`
	SubtitleEn *string `json:"subtitleEn"`
	SubtitlePl *string `json:"subtitlePl"`
	Slug       string  `json:"slug"`

	DescriptionEn *string `json:"descriptionEn"`
	DescriptionPl *string `json:"descriptionPl"`

	// Pricing
	Price    *decimal.Decimal `json:"price"`
	Currency *string          `json:"currency"`

	// HostIds IDs of hosts associated with the event.
	HostIds []uuid.UUID `json:"hostIds"`

	StartsAt time.Time `json:"startsAt"`
	EndsAt   time.Time `json:"endsAt"`

	IsVirtual        bool    `json:"isVirtual"`
	MeetingUrl       *string `json:"meetingUrl"`
	VenueNameEn      *string `json:"venueNameEn"`
	VenueNamePl      *string `json:"venueNamePl"`
	VenueStreet      *string `json:"venueStreet"`
	VenueCityEn      *string `json:"venueCityEn"`
	VenueCityPl      *string `json:"venueCityPl"`
	VenuePostalCode  *string `json:"venuePostalCode"`
	VenueCountryCode *string `json:"venueCountryCode"`
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

		validation.Field(&p.MeetingUrl, is.URL),

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

// HasEnded reports whether the event is over, in which case there is nothing
// left to register for. Recordings of paid events remain on sale.
func (d *EventDetailsDto) HasEnded() bool {
	return d.EndsAt.Before(time.Now().UTC())
}

type PublishEventValidation struct {
	DescriptionPl *string    `json:"descriptionPl"`
	DescriptionEn *string    `json:"descriptionEn"`
	PublishedAt   *time.Time `json:"publishedAt"`
}

func (p *PublishEventValidation) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.DescriptionEn, validation.Required),
		validation.Field(&p.DescriptionPl, validation.Required),
		validation.Field(&p.PublishedAt, validation.Nil.Error("event is already published")),
	)
}
