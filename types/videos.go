package types

import (
	"time"

	"github.com/moroz/homeosapiens-go/db/queries"
)

type VideoGroupListDTO struct {
	*queries.VideoGroup
	HasAccess     bool
	MinRecordedOn *time.Time
	MaxRecordedOn *time.Time
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

type VideoListDTO struct {
	*queries.Video
	Hosts []*queries.Host
}

func (v *VideoGroupListDTO) IsPremium() bool {
	return v.VideoGroup.ProductID != nil
}

func (v *VideoGroupDetailsDTO) IsPremium() bool {
	return v.VideoGroup.ProductID != nil
}
