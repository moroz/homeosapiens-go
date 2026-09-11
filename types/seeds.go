package types

import (
	"uuid"
)

type CreateVideoGroupParams struct {
	ID        uuid.UUID
	TitleEn   string
	TitlePl   string
	Slug      string
	ProductID *uuid.UUID
}

type CreateVideoParams struct {
	ID      uuid.UUID
	TitleEn string
	TitlePl string
	Slug    string
}
