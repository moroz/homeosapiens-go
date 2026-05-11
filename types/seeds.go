package types

import (
	"github.com/google/uuid"
)

type CreateVideoGroupParams struct {
	ID        uuid.UUID
	TitleEn   string
	TitlePl   string
	Slug      string
	ProductID *uuid.UUID
}
