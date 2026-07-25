package api

import (
	"context"
	"database/sql"
	"errors"

	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
)

type eventServer struct {
	db queries.DBTX
}

func NewEventServer(db queries.DBTX) *eventServer {
	return &eventServer{db: db}
}

func (s *eventServer) UpdateEvent(ctx context.Context, _ UpdateEventRequestObject) (UpdateEventResponseObject, error) {
	panic("Unimplemented")
}

func (s *eventServer) ListEvents(ctx context.Context, params ListEventsRequestObject) (ListEventsResponseObject, error) {
	page, perPage := resolvePaginationParams(params.Params.Page, params.Params.PerPage)

	events, err := queries.New(s.db).PaginateEvents(ctx, &queries.PaginateEventsParams{
		Page:    page,
		PerPage: perPage,
	})
	if err != nil {
		return nil, err
	}

	count, err := queries.New(s.db).CountEvents(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]Event, len(events))
	for i, e := range events {
		out[i] = Event{
			Id:         e.ID,
			Slug:       e.Slug,
			TitleEn:    e.TitleEn,
			TitlePl:    e.TitlePl,
			SubtitleEn: e.SubtitleEn,
			SubtitlePl: e.SubtitlePl,
			EventType:  string(e.EventType),
			IsVirtual:  e.IsVirtual,
			StartsAt:   e.StartsAt,
			EndsAt:     e.EndsAt,
			InsertedAt: e.InsertedAt,
			UpdatedAt:  e.UpdatedAt,
		}
	}

	return ListEvents200JSONResponse{
		Data: out,
		Pagination: Pagination{
			Page:       page,
			PerPage:    perPage,
			Total:      count,
			TotalPages: countPages(count, perPage),
		},
	}, nil
}

func (s *eventServer) CreateEvent(ctx context.Context, request CreateEventRequestObject) (CreateEventResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (s *eventServer) GetEvent(ctx context.Context, request GetEventRequestObject) (GetEventResponseObject, error) {
	e, err := services.NewEventService(s.db).GetEventDetailsById(ctx, request.Id, nil)
	if errors.Is(err, sql.ErrNoRows) {
		return GetEvent404Response{}, nil
	}
	if err != nil {
		return nil, err
	}

	price := new(string)
	currency := new(string)
	if !e.IsFree() {
		*price = e.Product.BasePriceAmount.StringFixedBank(2)
		*currency = e.Product.BasePriceCurrency
	}

	return GetEvent200JSONResponse{
		Currency:      currency,
		DescriptionEn: e.DescriptionEn,
		DescriptionPl: e.DescriptionPl,
		EndsAt:        e.EndsAt,
		EventType:     string(e.EventType),
		Hosts:         nil,
		Id:            e.ID,
		InsertedAt:    e.InsertedAt,
		IsFree:        e.IsFree(),
		IsVirtual:     e.IsVirtual,
		Price:         price,
		Slug:          e.Slug,
		StartsAt:      e.StartsAt,
		SubtitleEn:    e.SubtitleEn,
		SubtitlePl:    e.SubtitlePl,
		TitleEn:       e.TitleEn,
		TitlePl:       e.TitlePl,
		UpdatedAt:     e.UpdatedAt,
	}, nil
}
