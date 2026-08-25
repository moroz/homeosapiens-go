package api

import (
	"context"
	"database/sql"
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/types"
)

type videoServer struct {
	db queries.DBTX
	q  *queries.Queries
}

func NewVideoServer(db queries.DBTX) *videoServer {
	return &videoServer{db: db, q: queries.New(db)}
}

// videoDetails maps a video row onto the VideoDetails schema, the payload the
// admin edit form reads and writes.
func videoDetails(v *queries.Video) VideoDetails {
	return VideoDetails{
		Id:              v.ID,
		Slug:            v.Slug,
		TitleEn:         v.TitleEn,
		TitlePl:         v.TitlePl,
		DescriptionEn:   v.DescriptionEn,
		DescriptionPl:   v.DescriptionPl,
		Provider:        v.Provider,
		IsPublic:        v.IsPublic,
		RecordedOn:      v.RecordedOn,
		HostId:          v.HostID,
		YoutubeId:       v.YoutubeID,
		DurationSeconds: v.DurationSeconds,
		InsertedAt:      v.InsertedAt,
		UpdatedAt:       v.UpdatedAt,
	}
}

func (s *videoServer) ListVideos(ctx context.Context, params ListVideosRequestObject) (ListVideosResponseObject, error) {
	page, perPage := resolvePaginationParams(params.Params.Page, params.Params.PerPage)

	videos, err := s.q.PaginateVideos(ctx, &queries.PaginateVideosParams{
		Page:    page,
		PerPage: perPage,
	})
	if err != nil {
		return nil, err
	}

	count, err := s.q.CountVideos(ctx)
	if err != nil {
		return nil, err
	}

	return ListVideos200JSONResponse{
		Data: videoList(videos),
		Pagination: Pagination{
			Page:       page,
			PerPage:    perPage,
			Total:      count,
			TotalPages: countPages(count, perPage),
		},
	}, nil
}

func (s *videoServer) GetVideo(ctx context.Context, request GetVideoRequestObject) (GetVideoResponseObject, error) {
	v, err := s.q.GetVideoById(ctx, request.Id)
	if errors.Is(err, sql.ErrNoRows) {
		return GetVideo404Response{}, nil
	}
	if err != nil {
		return nil, err
	}

	return GetVideo200JSONResponse(videoDetails(v)), nil
}

func (s *videoServer) UpdateVideo(ctx context.Context, request UpdateVideoRequestObject) (UpdateVideoResponseObject, error) {
	p := request.Body

	v, err := services.NewVideoService(s.db).UpdateVideo(ctx, request.Id, &types.UpdateVideoInput{
		TitleEn:       p.TitleEn,
		TitlePl:       p.TitlePl,
		Slug:          p.Slug,
		DescriptionEn: p.DescriptionEn,
		DescriptionPl: p.DescriptionPl,
		RecordedOn:    p.RecordedOn,
		IsPublic:      p.IsPublic,
		HostID:        p.HostId,
	})
	if errors.Is(err, sql.ErrNoRows) {
		return UpdateVideo404Response{}, nil
	}
	if verr, ok := errors.AsType[validation.Errors](err); ok {
		return UpdateVideo422JSONResponse{Errors: validationErrorMessages(verr)}, nil
	}
	if err != nil {
		return nil, err
	}

	return UpdateVideo200JSONResponse(videoDetails(v)), nil
}
