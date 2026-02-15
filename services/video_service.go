package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/moroz/homeosapiens-go/db/queries"
)

type VideoService struct {
	db queries.DBTX
}

func NewVideoService(db queries.DBTX) *VideoService {
	return &VideoService{db}
}

type VideoListDto struct {
	*queries.Video
	Sources []*queries.VideoSource
}

func (s *VideoService) ListVideosWithSources(ctx context.Context, user *queries.User) ([]*VideoListDto, error) {
	var videos []*queries.Video
	var err error

	if user == nil {
		videos, err = queries.New(s.db).ListPublicVideos(ctx)
	} else {
		videos, err = queries.New(s.db).ListVideosForUser(ctx, user.ID)
	}
	if err != nil {
		return nil, err
	}

	var ids []uuid.UUID
	for _, v := range videos {
		ids = append(ids, v.ID)
	}

	sources, err := s.preloadSourcesForVideos(ctx, ids)
	if err != nil {
		return nil, err
	}

	var result []*VideoListDto
	for _, v := range videos {
		result = append(result, &VideoListDto{
			Video:   v,
			Sources: sources[v.ID],
		})
	}

	return result, nil
}

func (s *VideoService) preloadSourcesForVideos(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID][]*queries.VideoSource, error) {
	sources, err := queries.New(s.db).ListVideoSourcesForVideos(ctx, ids)
	if err != nil {
		return nil, err
	}

	var result = make(map[uuid.UUID][]*queries.VideoSource)
	for _, row := range sources {
		result[row.VideoID] = append(result[row.VideoID], row)
	}
	return result, nil
}
