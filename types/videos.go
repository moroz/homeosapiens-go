package types

import "github.com/moroz/homeosapiens-go/db/queries"

type VideoGroupListDTO struct {
	*queries.VideoGroup
}

type VideoListDTO struct {
	*queries.Video
	Thumbnail *queries.Asset
}

type VideoGroupDetailsDTO struct {
	*queries.VideoGroup
	Videos []*VideoListDTO
}

type VideoDetailsDTO struct {
	*queries.Video
	Sources []*queries.VideoSource
}
