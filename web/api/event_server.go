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

type eventServer struct {
	db queries.DBTX
}

func NewEventServer(db queries.DBTX) *eventServer {
	return &eventServer{db: db}
}

func (s *eventServer) UpdateEvent(ctx context.Context, request UpdateEventRequestObject) (UpdateEventResponseObject, error) {
	svc := services.NewEventService(s.db)

	if _, err := svc.UpdateEvent(ctx, request.Id, request.Body); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return UpdateEvent404Response{}, nil
		}
		if verr, ok := errors.AsType[validation.Errors](err); ok {
			return UpdateEvent422JSONResponse{Errors: validationErrorMessages(verr)}, nil
		}
		return nil, err
	}

	e, err := svc.GetEventDetailsById(ctx, request.Id, nil)
	if err != nil {
		return nil, err
	}

	price := new(string)
	currency := new(string)
	if !e.IsFree() {
		*price = e.Product.BasePriceAmount.StringFixedBank(2)
		*currency = e.Product.BasePriceCurrency
	}

	return UpdateEvent200JSONResponse{
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
		Currency:         currency,
		DescriptionEn:    e.DescriptionEn,
		DescriptionPl:    e.DescriptionPl,
		EndsAt:           e.EndsAt,
		EventType:        string(e.EventType),
		Hosts:            nil,
		Id:               e.ID,
		InsertedAt:       e.InsertedAt,
		IsFree:           e.IsFree(),
		IsVirtual:        e.IsVirtual,
		Price:            price,
		Slug:             e.Slug,
		StartsAt:         e.StartsAt,
		SubtitleEn:       e.SubtitleEn,
		SubtitlePl:       e.SubtitlePl,
		TitleEn:          e.TitleEn,
		TitlePl:          e.TitlePl,
		UpdatedAt:        e.UpdatedAt,
		VenueCityEn:      e.VenueCityEn,
		VenueCityPl:      e.VenueCityPl,
		VenueCountryCode: e.VenueCountryCode,
		VenueNameEn:      e.VenueNameEn,
		VenueNamePl:      e.VenueNamePl,
		VenuePostalCode:  e.VenuePostalCode,
		VenueStreet:      e.VenueStreet,
	}, nil
}

func (s *eventServer) CreateEvent(ctx context.Context, request CreateEventRequestObject) (CreateEventResponseObject, error) {
	p := request.Body

	var price *decimal.Decimal
	if p.Price != nil {
		parsed, err := decimal.NewFromString(*p.Price)
		if err != nil {
			return nil, err
		}
		price = &parsed
	}

	input := &types.CreateEventInput{
		EventType:        p.EventType,
		TitleEn:          p.TitleEn,
		TitlePl:          p.TitlePl,
		SubtitleEn:       p.SubtitleEn,
		SubtitlePl:       p.SubtitlePl,
		Slug:             p.Slug,
		DescriptionEn:    p.DescriptionEn,
		DescriptionPl:    p.DescriptionPl,
		Price:            price,
		Currency:         p.Currency,
		HostIds:          p.HostIds,
		StartsAt:         p.StartsAt,
		EndsAt:           p.EndsAt,
		IsVirtual:        p.IsVirtual,
		VenueNameEn:      p.VenueNameEn,
		VenueNamePl:      p.VenueNamePl,
		VenueStreet:      p.VenueStreet,
		VenueCityEn:      p.VenueCityEn,
		VenueCityPl:      p.VenueCityPl,
		VenuePostalCode:  p.VenuePostalCode,
		VenueCountryCode: p.VenueCountryCode,
	}

	e, err := services.NewEventService(s.db).CreateEvent(ctx, input)
	if err != nil {
		if err, ok := errors.AsType[validation.Errors](err); ok {
			return CreateEvent422JSONResponse{Errors: validationErrorMessages(err)}, nil
		}
		return nil, err
	}

	eventPrice := new(string)
	currency := new(string)
	if !e.IsFree() {
		*eventPrice = e.Product.BasePriceAmount.StringFixedBank(2)
		*currency = e.Product.BasePriceCurrency
	}

	location := fmt.Sprintf("/events/%s", e.ID)

	return CreateEvent201JSONResponse{
		Headers: CreateEvent201ResponseHeaders{
			Location: &location,
		},
		Body: EventDetails{
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
			Price:         eventPrice,
			Slug:          e.Slug,
			StartsAt:      e.StartsAt,
			SubtitleEn:    e.SubtitleEn,
			SubtitlePl:    e.SubtitlePl,
			TitleEn:       e.TitleEn,
			TitlePl:       e.TitlePl,
			UpdatedAt:     e.UpdatedAt,
		},
	}, nil
}

func validationErrorMessages(verrs validation.Errors) map[string]string {
	messages := make(map[string]string, len(verrs))
	for field, ferr := range verrs {
		messages[field] = ferr.Error()
	}
	return messages
}
