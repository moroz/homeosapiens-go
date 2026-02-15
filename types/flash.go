package types

import (
	"encoding/gob"

	"github.com/google/uuid"
)

func init() {
	gob.Register(make(Flash))
	gob.Register(uuid.UUID{})
}

type Flash map[string]string
