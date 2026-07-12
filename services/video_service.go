package services

import (
	"cmp"
	"context"
	"slices"

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

func (s *VideoService) ListVideoGroupsForUser(ctx context.Context, userID uuid.UUID) ([]*types.VideoGroupListDTO, error) {
	videos, err := queries.New(s.db).ListVideoGroupsForUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	ids := make([]uuid.UUID, len(videos))
	for i, v := range videos {
		ids[i] = v.VideoGroup.ID
	}

	dateRanges, err := queries.New(s.db).GetMinMaxRecordedDatesForVideoGroups(ctx, ids)
	if err != nil {
		return nil, err
	}

	dateRangeMap := make(map[uuid.UUID]*queries.GetMinMaxRecordedDatesForVideoGroupsRow)
	for _, row := range dateRanges {
		dateRangeMap[row.ID] = row
	}

	var result []*types.VideoGroupListDTO
	for _, vg := range videos {

		dto := &types.VideoGroupListDTO{
			VideoGroup: &vg.VideoGroup,
			HasAccess:  vg.HasAccess,
		}

		if dateRange, ok := dateRangeMap[vg.VideoGroup.ID]; ok {
			dto.MinRecordedOn = &dateRange.MinRecordedOn
			dto.MaxRecordedOn = &dateRange.MaxRecordedOn
		}

		result = append(result, dto)
	}

	// Sort the video groups in go, descending by the lowest video recorded on date. Consider doing this in SQL if this
	// syntax becomes too cumbersome (might require CTEs).
	result = slices.SortedFunc(slices.Values(result), func(a *types.VideoGroupListDTO, b *types.VideoGroupListDTO) int {
		var aUnix, bUnix int64
		if a.MinRecordedOn != nil {
			aUnix = a.MinRecordedOn.Unix()
		}
		if b.MinRecordedOn != nil {
			bUnix = b.MinRecordedOn.Unix()
		}

		// Compare in reverse order to sort descending
		return cmp.Compare(bUnix, aUnix)
	})

	return result, nil
}

func (s *VideoService) GetVideoGroupDetails(ctx context.Context, userID uuid.UUID, slug *string) (*types.VideoGroupDetailsDTO, error) {
	group, err := queries.New(s.db).GetVideoGroupForUserBySlug(ctx, &queries.GetVideoGroupForUserBySlugParams{
		Slug:   slug,
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	videos, err := queries.New(s.db).ListVideosForVideoGroup(ctx, group.VideoGroup.ID)
	if err != nil {
		return nil, err
	}

	return &types.VideoGroupDetailsDTO{
		HasAccess:  group.HasAccess,
		VideoGroup: &group.VideoGroup,
		Videos:     videos,
	}, nil
}

func (s *VideoService) GetVideoThumbnailData(ctx context.Context, id uuid.UUID) (*queries.GetVideoThumbnailDataRow, error) {
	return queries.New(s.db).GetVideoThumbnailData(ctx, id)
}

func (s *VideoService) GetVideoForUser(ctx context.Context, userID uuid.UUID, groupSlug, videoSlug string) (*types.VideoDetailsDTO, error) {
	video, err := queries.New(s.db).GetVideoForUser(ctx, &queries.GetVideoForUserParams{
		VideoSlug: videoSlug,
		GroupSlug: groupSlug,
		UserID:    userID,
	})
	if err != nil {
		return nil, err
	}

	result := &types.VideoDetailsDTO{
		Video:     &video.Video,
		HasAccess: video.HasAccess,
	}

	if video.HasAccess {
		sources, err := queries.New(s.db).ListVideoSourcesForVideos(ctx, []uuid.UUID{video.Video.ID})
		if err != nil {
			return nil, err
		}

		result.Sources = sources
	}

	return result, nil
}
