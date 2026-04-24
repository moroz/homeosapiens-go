package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/types"
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

func (s *VideoService) CreateVideoGroup(ctx context.Context, params *types.CreateVideoGroupParams) (*queries.VideoGroup, error) {
	return queries.New(s.db).InsertVideoGroup(ctx, &queries.InsertVideoGroupParams{
		TitleEn:   params.TitleEn,
		TitlePl:   params.TitlePl,
		Slug:      params.Slug,
		ProductID: params.ProductID,
	})
}

func (s *VideoService) ListVideoGroupsForUser(ctx context.Context, user *queries.User) ([]*types.VideoGroupListDTO, error) {
	var userID *uuid.UUID
	if user != nil {
		userID = &user.ID
	}

	videos, err := queries.New(s.db).ListVideoGroupsForUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	var result []*types.VideoGroupListDTO
	for _, vg := range videos {
		result = append(result, &types.VideoGroupListDTO{
			VideoGroup: vg,
		})
	}

	return result, nil
}

func (s *VideoService) GetVideoGroupDetails(ctx context.Context, user *queries.User, slug *string) (*types.VideoGroupDetailsDTO, error) {
	var userID *uuid.UUID
	if user != nil {
		userID = &user.ID
	}

	group, err := queries.New(s.db).GetVideoGroupForUserBySlug(ctx, &queries.GetVideoGroupForUserBySlugParams{
		Slug:   slug,
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	videos, err := queries.New(s.db).ListVideosForVideoGroup(ctx, group.ID)
	if err != nil {
		return nil, err
	}

	return &types.VideoGroupDetailsDTO{
		VideoGroup: group,
		Videos:     videos,
	}, nil
}

func (s *VideoService) GetVideoForUser(ctx context.Context, user *queries.User, groupSlug, videoSlug string) (*types.VideoDetailsDTO, error) {
	var userID *uuid.UUID
	if user != nil {
		userID = &user.ID
	}

	video, err := queries.New(s.db).GetVideoForUser(ctx, &queries.GetVideoForUserParams{
		VideoSlug: videoSlug,
		GroupSlug: groupSlug,
		UserID:    userID,
	})
	if err != nil {
		return nil, err
	}

	sources, err := queries.New(s.db).ListVideoSourcesForVideos(ctx, []uuid.UUID{video.ID})
	if err != nil {
		return nil, err
	}

	return &types.VideoDetailsDTO{
		Video:   video,
		Sources: sources,
	}, nil
}
