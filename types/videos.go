package types

import "github.com/moroz/homeosapiens-go/db/queries"

type VideoGroupListDTO struct {
	*queries.VideoGroup
}

type VideoGroupDetailsDTO struct {
	*queries.VideoGroup
	Videos []*queries.ListVideosForVideoGroupRow
}

type VideoDetailsDTO struct {
	*queries.Video
	Sources []*queries.VideoSource
}
