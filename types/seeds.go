package types

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/moroz/homeosapiens-go/db/queries"
)

type CreateAssetParams struct {
	ID               pgtype.UUID
	ObjectKey        string
	OriginalFilename string
}

type CreateHostParams struct {
	ID               pgtype.UUID
	Salutation       string
	GivenName        string
	FamilyName       string
	ProfilePictureId pgtype.UUID
	Country          string
}

type CreateVenueParams struct {
	ID          pgtype.UUID
	NameEn      string
	NamePl      *string
	CityEn      string
	CityPl      *string
	CountryCode string
	Street      string
	PostalCode  string
}

type CreateEventParams struct {
	ID                pgtype.UUID
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
	VenueID           pgtype.UUID
	BasePriceAmount   *string
	BasePriceCurrency *string
}

type CreateEventHostParams struct {
	EventID  pgtype.UUID
	HostID   pgtype.UUID
	Position int32
}

type CreateEventPriceParams struct {
	ID            pgtype.UUID
	EventID       pgtype.UUID
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
	ID       pgtype.UUID
	EventID  pgtype.UUID
	Provider queries.VideoProvider
	TitleEn  string
	TitlePl  string
	Slug     string
	IsPublic bool
}

type CreateVideoSourceParams struct {
	ID          pgtype.UUID
	VideoID     pgtype.UUID
	ContentType string
	Codec       string
	ObjectKey   string
}
