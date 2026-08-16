package types

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	GivenName  string    `json:"givenName"`
	FamilyName string    `json:"familyName"`
	Role       string    `json:"role"`
	InsertedAt time.Time `json:"insertedAt"`
}
