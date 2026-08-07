package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/types"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/shopspring/decimal"
)

type eventServer struct {
	db *pgxpool.Pool
}

func NewEventServer(db *pgxpool.Pool) *eventServer {
	return &eventServer{db: db}
}

// eventDetails maps an event and its associations onto the EventDetails schema.
// Every operation returning event details goes through it, so that a client can
// treat the payloads of GET, POST and PATCH interchangeably.
func eventDetails(e *types.EventDetailsDto) EventDetails {
	// price and currency stay null for a free event, as documented in the schema.
	var price, currency *string
	if !e.IsFree() {
		price = new(e.Product.BasePriceAmount.StringFixedBank(2))
		currency = new(e.Product.BasePriceCurrency)
	}

	hosts := make([]Host, len(e.Hosts))
	for i, host := range e.Hosts {
		hosts[i] = Host{
			Country:    host.Country,
			FamilyName: host.FamilyName,
			GivenName:  host.GivenName,
			Id:         host.ID,
			Salutation: host.Salutation,
		}
	}

	return EventDetails{
		Currency:         currency,
		DescriptionEn:    e.DescriptionEn,
		DescriptionPl:    e.DescriptionPl,
		EndsAt:           e.EndsAt,
		EventType:        string(e.EventType),
		Hosts:            hosts,
		Id:               e.ID,
		InsertedAt:       e.InsertedAt,
		IsFree:           e.IsFree(),
		IsVirtual:        e.IsVirtual,
		MeetingUrl:       e.MeetingUrl,
		Price:            price,
		PublishedAt:      e.PublishedAt,
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
	}
}

func (s *eventServer) UpdateEvent(ctx context.Context, request UpdateEventRequestObject) (UpdateEventResponseObject, error) {
	e, err := services.NewEventService(s.db).UpdateEvent(ctx, request.Id, request.Body)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return UpdateEvent404Response{}, nil
		}
		if verr, ok := errors.AsType[validation.Errors](err); ok {
			return UpdateEvent422JSONResponse{Errors: validationErrorMessages(verr)}, nil
		}
		return nil, err
	}

	return UpdateEvent200JSONResponse(eventDetails(e)), nil
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
			Id:          e.ID,
			Slug:        e.Slug,
			TitleEn:     e.TitleEn,
			TitlePl:     e.TitlePl,
			SubtitleEn:  e.SubtitleEn,
			SubtitlePl:  e.SubtitlePl,
			EventType:   string(e.EventType),
			IsVirtual:   e.IsVirtual,
			MeetingUrl:  e.MeetingUrl,
			StartsAt:    e.StartsAt,
			EndsAt:      e.EndsAt,
			InsertedAt:  e.InsertedAt,
			UpdatedAt:   e.UpdatedAt,
			PublishedAt: e.PublishedAt,
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

	return GetEvent200JSONResponse(eventDetails(e)), nil
}

func (s *eventServer) CreateEvent(ctx context.Context, request CreateEventRequestObject) (CreateEventResponseObject, error) {
	p := request.Body

	var price *decimal.Decimal
	if p.Price != nil {
		parsed, err := decimal.NewFromString(*p.Price)
		if err != nil {
			return CreateEvent422JSONResponse{Errors: map[string]string{
				"price": "must be a decimal number",
			}}, nil
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
		MeetingUrl:       p.MeetingUrl,
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

	location := fmt.Sprintf("/events/%s", e.ID)

	return CreateEvent201JSONResponse{
		Headers: CreateEvent201ResponseHeaders{
			Location: &location,
		},
		Body: eventDetails(e),
	}, nil
}

func (s *eventServer) PublishEvent(ctx context.Context, request PublishEventRequestObject) (PublishEventResponseObject, error) {
	_, err := services.NewEventService(s.db).PublishEvent(ctx, request.Id)
	if errors.Is(err, sql.ErrNoRows) {
		return PublishEvent404Response{}, nil
	}
	if err, ok := errors.AsType[validation.Errors](err); ok {
		return PublishEvent422JSONResponse{Errors: validationErrorMessages(err)}, nil
	}
	if err != nil {
		return nil, err
	}

	return PublishEvent204Response{}, nil
}

func (s *eventServer) UnpublishEvent(ctx context.Context, request UnpublishEventRequestObject) (UnpublishEventResponseObject, error) {
	_, err := services.NewEventService(s.db).UnpublishEvent(ctx, request.Id)
	if errors.Is(err, sql.ErrNoRows) {
		return UnpublishEvent404Response{}, nil
	}
	if err != nil {
		return nil, err
	}

	return UnpublishEvent204Response{}, nil
}

func (s *eventServer) DeleteEvent(ctx context.Context, request DeleteEventRequestObject) (DeleteEventResponseObject, error) {
	err := services.NewEventService(s.db).DeleteEvent(ctx, request.Id)
	if errors.Is(err, sql.ErrNoRows) {
		return DeleteEvent404Response{}, nil
	}
	if errors.Is(err, services.ErrEventHasRegistrations) {
		return DeleteEvent409JSONResponse{Errors: map[string]string{
			"registrations": "event has registrations and cannot be deleted; unpublish it instead",
		}}, nil
	}
	if err != nil {
		return nil, err
	}

	return DeleteEvent204Response{}, nil
}

// ListEventAttendants answers the attendants endpoint as currently specified: the
// operation declares no response body yet, so this only distinguishes a known
// event from an unknown one. The payload is still to be defined in openapi.yaml.
func (s *eventServer) ListEventAttendants(ctx context.Context, request ListEventAttendantsRequestObject) (ListEventAttendantsResponseObject, error) {
	page, perPage := resolvePaginationParams(request.Params.Page, request.Params.PerPage)

	list, err := services.NewEventRegistrationService(s.db).PaginateEventAttendants(ctx, &queries.PaginateEventRegistrationsParams{
		EventID: request.Id,
		Page:    page,
		PerPage: perPage,
	})
	if errors.Is(err, sql.ErrNoRows) {
		return ListEventAttendants404Response{}, nil
	}
	if err != nil {
		return nil, err
	}

	result := make([]EventAttendant, len(list.Data))
	for i, row := range list.Data {
		var orderNumber *string
		if row.OrderNumber != nil {
			orderNumber = new(fmt.Sprintf("%d", *row.OrderNumber))
		}

		result[i] = EventAttendant{
			Id:          row.ID,
			Email:       openapi_types.Email(row.Email.Plaintext()),
			FamilyName:  row.FamilyName.Plaintext(),
			GivenName:   row.GivenName.Plaintext(),
			OrderId:     row.OrderID,
			OrderNumber: orderNumber,
			InsertedAt:  row.InsertedAt,
		}
	}

	return ListEventAttendants200JSONResponse{
		Data: result,
		Pagination: Pagination{
			Page:       int32(list.Pagination.Page),
			PerPage:    int32(list.Pagination.PerPage),
			TotalPages: int32(list.Pagination.TotalPages),
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
