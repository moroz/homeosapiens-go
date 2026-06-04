package types

import "github.com/moroz/homeosapiens-go/db/queries"

type VideoGroupListDTO struct {
	*queries.VideoGroup
	HasAccess bool
}

type VideoGroupDetailsDTO struct {
	*queries.VideoGroup
	HasAccess bool
	Videos    []*queries.Video
}

type VideoDetailsDTO struct {
	*queries.Video
	HasAccess bool
	Sources   []*queries.VideoSource
}
