package api

import (
	"math"

	"github.com/moroz/homeosapiens-go/config"
)

type AsInt32 interface {
	~int32
}

func resolvePaginationParams[P AsInt32](pageParam *P, perPageParam *P) (page int32, perPage int32) {
	page = 1
	perPage = config.DefaultPageSize

	if pageParam != nil && *pageParam > 0 {
		page = int32(*pageParam)
	}

	if perPageParam != nil && *perPageParam > 0 {
		perPage = int32(*perPageParam)
	}

	return
}

func countPages(count int64, perPage int32) int32 {
	return int32(
		math.Ceil(
			float64(count) / float64(perPage),
		),
	)
}
