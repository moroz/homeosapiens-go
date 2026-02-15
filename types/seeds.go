package types

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/moroz/homeosapiens-go/db/queries"
)

type CreateAssetParams struct {
	ID               uuid.UUID
	ObjectKey        string
	OriginalFilename string
}

type CreateHostParams struct {
	ID               uuid.UUID
	Salutation       string
	GivenName        string
	FamilyName       string
	ProfilePictureId uuid.UUID
	Country          string
}

type CreateVenueParams struct {
	ID          uuid.UUID
	NameEn      string
	NamePl      *string
	CityEn      string
	CityPl      *string
	CountryCode string
	Street      string
	PostalCode  string
}

type CreateEventParams struct {
	ID                uuid.UUID
	EventType         queries.EventType
	TitleEn           string
	TitlePl           string
	SubtitleEn        *string
	SubtitlePl        *string
	Slug              string
	StartsAt          pgtype.Timestamp
	EndsAt            pgtype.Timestamp
	IsVirtual         bool
	DescriptionEn     string
	DescriptionPl     *string
	BasePriceAmount   *string
	BasePriceCurrency *string
	VenueNameEn       *string
	VenueNamePl       *string
	VenueCityEn       *string
	VenueCityPl       *string
	VenueStreet       *string
	VenuePostalCode   *string
	VenueCountryCode  *string
}

type CreateEventHostParams struct {
	EventID  uuid.UUID
	HostID   uuid.UUID
	Position int32
}

type CreateEventPriceParams struct {
	ID            uuid.UUID
	EventID       uuid.UUID
	PriceType     queries.PriceType
	RuleType      queries.PriceRuleType
	PriceAmount   string
	PriceCurrency string
	DiscountCode  *string
	Priority      int32
	IsActive      bool
	ValidFrom     pgtype.Timestamp
	ValidUntil    pgtype.Timestamp
}

type CreateVideoParams struct {
	ID       uuid.UUID
	EventID  uuid.UUID
	Provider queries.VideoProvider
	TitleEn  string
	TitlePl  string
	Slug     string
	IsPublic bool
}

type CreateVideoSourceParams struct {
	ID          uuid.UUID
	VideoID     uuid.UUID
	ContentType string
	Codec       string
	ObjectKey   string
}
