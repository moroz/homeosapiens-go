package types

import "encoding/gob"

func init() {
	gob.Register(make(Flash))
}

type Flash map[string]string
