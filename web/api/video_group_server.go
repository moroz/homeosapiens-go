package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/shopspring/decimal"
)

type videoGroupServer struct {
	db queries.DBTX
}

func NewVideoGroupServer(db queries.DBTX) *videoGroupServer {
	return &videoGroupServer{db: db}
}

// videoGroup maps a group onto the VideoGroup schema. Every operation returning
// a group goes through it, so that a client can treat the payloads of GET, POST
// and PATCH interchangeably.
func videoGroup(g *services.VideoGroupAdminDto) VideoGroup {
	// price and currency stay null for a free group, as documented in the schema.
	var price, currency *string
	if g.IsPremium() && g.Price != nil {
		price = new(g.Price.StringFixedBank(2))
		currency = g.Currency
	}

	return VideoGroup{
		Currency:   currency,
		Id:         g.ID,
		InsertedAt: g.InsertedAt,
		IsPremium:  g.IsPremium(),
		Price:      price,
		Slug:       g.Slug,
		TitleEn:    g.TitleEn,
		TitlePl:    g.TitlePl,
		UpdatedAt:  g.UpdatedAt,
		VideoCount: int32(g.VideoCount),
	}
}

func (s *videoGroupServer) ListVideoGroups(ctx context.Context, request ListVideoGroupsRequestObject) (ListVideoGroupsResponseObject, error) {
	page, perPage := resolvePaginationParams(request.Params.Page, request.Params.PerPage)

	groups, count, err := services.NewVideoGroupService(s.db).ListVideoGroups(ctx, page, perPage)
	if err != nil {
		return nil, err
	}

	out := make([]VideoGroup, len(groups))
	for i, g := range groups {
		out[i] = videoGroup(g)
	}

	return ListVideoGroups200JSONResponse{
		Data: out,
		Pagination: Pagination{
			Page:       page,
			PerPage:    perPage,
			Total:      count,
			TotalPages: countPages(count, perPage),
		},
	}, nil
}

func (s *videoGroupServer) GetVideoGroup(ctx context.Context, request GetVideoGroupRequestObject) (GetVideoGroupResponseObject, error) {
	g, err := services.NewVideoGroupService(s.db).GetVideoGroupById(ctx, request.Id)
	if errors.Is(err, sql.ErrNoRows) {
		return GetVideoGroup404Response{}, nil
	}
	if err != nil {
		return nil, err
	}

	return GetVideoGroup200JSONResponse(videoGroup(g)), nil
}

func (s *videoGroupServer) CreateVideoGroup(ctx context.Context, request CreateVideoGroupRequestObject) (CreateVideoGroupResponseObject, error) {
	p := request.Body

	var price *decimal.Decimal
	if p.Price != nil {
		parsed, err := decimal.NewFromString(*p.Price)
		if err != nil {
			return CreateVideoGroup422JSONResponse{Errors: map[string]string{
				"price": "must be a decimal number",
			}}, nil
		}
		price = &parsed
	}

	g, err := services.NewVideoGroupService(s.db).CreateVideoGroup(ctx, &types.CreateVideoGroupInput{
		TitleEn:  p.TitleEn,
		TitlePl:  p.TitlePl,
		Slug:     p.Slug,
		Price:    price,
		Currency: p.Currency,
	})
	if verr, ok := errors.AsType[validation.Errors](err); ok {
		return CreateVideoGroup422JSONResponse{Errors: validationErrorMessages(verr)}, nil
	}
	if err != nil {
		return nil, err
	}

	location := fmt.Sprintf("/api/admin/video-groups/%s", g.ID)

	return CreateVideoGroup201JSONResponse{
		Headers: CreateVideoGroup201ResponseHeaders{
			Location: &location,
		},
		Body: videoGroup(g),
	}, nil
}

func (s *videoGroupServer) ListVideosInVideoGroup(ctx context.Context, request ListVideosInVideoGroupRequestObject) (ListVideosInVideoGroupResponseObject, error) {
	videos, err := services.NewVideoGroupService(s.db).ListVideosInVideoGroup(ctx, request.Id)
	if errors.Is(err, sql.ErrNoRows) {
		return ListVideosInVideoGroup404Response{}, nil
	}
	if err != nil {
		return nil, err
	}

	return ListVideosInVideoGroup200JSONResponse(videoList(videos)), nil
}

func (s *videoGroupServer) ReplaceVideosInVideoGroup(ctx context.Context, request ReplaceVideosInVideoGroupRequestObject) (ReplaceVideosInVideoGroupResponseObject, error) {
	videos, err := services.NewVideoGroupService(s.db).ReplaceVideosInVideoGroup(ctx, request.Id, request.Body.VideoIds)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ReplaceVideosInVideoGroup404Response{}, nil
		}
		if verr, ok := errors.AsType[validation.Errors](err); ok {
			return ReplaceVideosInVideoGroup422JSONResponse{Errors: validationErrorMessages(verr)}, nil
		}
		return nil, err
	}

	return ReplaceVideosInVideoGroup200JSONResponse(videoList(videos)), nil
}

func videoList(videos []*queries.Video) []Video {
	out := make([]Video, len(videos))
	for i, v := range videos {
		out[i] = Video{
			Id:         v.ID,
			Slug:       v.Slug,
			TitleEn:    v.TitleEn,
			TitlePl:    v.TitlePl,
			Provider:   v.Provider,
			IsPublic:   v.IsPublic,
			RecordedOn: v.RecordedOn,
		}
	}
	return out
}

func (s *videoGroupServer) UpdateVideoGroup(ctx context.Context, request UpdateVideoGroupRequestObject) (UpdateVideoGroupResponseObject, error) {
	g, err := services.NewVideoGroupService(s.db).UpdateVideoGroup(ctx, request.Id, request.Body)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return UpdateVideoGroup404Response{}, nil
		}
		if verr, ok := errors.AsType[validation.Errors](err); ok {
			return UpdateVideoGroup422JSONResponse{Errors: validationErrorMessages(verr)}, nil
		}
		return nil, err
	}

	return UpdateVideoGroup200JSONResponse(videoGroup(g)), nil
}
