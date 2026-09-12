package types

import (
	"encoding/gob"

	"uuid"
)

func init() {
	gob.Register(make(Flash))
	gob.Register(uuid.UUID{})
}

type Flash map[string]string
