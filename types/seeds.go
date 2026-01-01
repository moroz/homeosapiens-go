package types

import (
	"github.com/jackc/pgx/v5/pgtype"
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
