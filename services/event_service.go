package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/moroz/homeosapiens-go/db/queries"
)

type EventService struct {
	db queries.DBTX
}

func NewEventService(db queries.DBTX) *EventService {
	return &EventService{db}
}

func (s *EventService) GetEventById(ctx context.Context, id string) (*queries.Event, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	return queries.New(s.db).GetEventById(ctx, pgtype.UUID{Valid: true, Bytes: uuid})
}

type EventListDto struct {
	*queries.ListEventsRow
	Hosts  []*queries.ListHostsForEventsRow
	Prices []*queries.EventPrice
}

func (s *EventService) ListEvents(ctx context.Context) ([]*EventListDto, error) {
	events, err := queries.New(s.db).ListEvents(ctx)
	if err != nil {
		return nil, err
	}

	var ids []pgtype.UUID
	for _, event := range events {
		ids = append(ids, event.ID)
	}

	hosts, err := s.preloadHostsForEvents(ctx, ids)
	if err != nil {
		return nil, err
	}

	prices, err := s.preloadPricesForEvents(ctx, ids)
	if err != nil {
		return nil, err
	}

	var result []*EventListDto
	for _, event := range events {
		result = append(result, &EventListDto{
			ListEventsRow: event,
			Hosts:         hosts[event.ID],
			Prices:        prices[event.ID],
		})
	}

	return result, nil
}

type EventDetailsDto struct {
	*queries.Event
	*queries.Venue
	Prices []*queries.EventPrice
	Hosts  []*queries.ListHostsForEventsRow
}

func (s *EventService) GetEventBySlug(ctx context.Context, slug string) (*EventDetailsDto, error) {
	var dto EventDetailsDto

	q := queries.New(s.db)
	event, err := q.GetEventBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	dto.Event = event

	if event.VenueID.Valid {
		venue, err := q.GetVenueById(ctx, event.VenueID)
		if err != nil {
			return nil, err
		}
		dto.Venue = venue
	}

	prices, err := s.preloadPricesForEvents(ctx, []pgtype.UUID{event.ID})
	if err != nil {
		return nil, err
	}
	dto.Prices = prices[event.ID]

	hosts, err := s.preloadHostsForEvents(ctx, []pgtype.UUID{event.ID})
	if err != nil {
		return nil, err
	}
	dto.Hosts = hosts[event.ID]

	return &dto, nil
}

func (s *EventService) preloadHostsForEvents(ctx context.Context, eventIds []pgtype.UUID) (map[pgtype.UUID][]*queries.ListHostsForEventsRow, error) {
	hosts, err := queries.New(s.db).ListHostsForEvents(ctx, eventIds)
	if err != nil {
		return nil, err
	}

	hostMap := make(map[pgtype.UUID][]*queries.ListHostsForEventsRow)
	for _, row := range hosts {
		hostMap[row.EventID] = append(hostMap[row.EventID], row)
	}

	return hostMap, nil
}

func (s *EventService) preloadPricesForEvents(ctx context.Context, eventIds []pgtype.UUID) (map[pgtype.UUID][]*queries.EventPrice, error) {
	prices, err := queries.New(s.db).ListPricesForEvents(ctx, eventIds)
	if err != nil {
		return nil, err
	}

	priceMap := make(map[pgtype.UUID][]*queries.EventPrice)
	for _, row := range prices {
		priceMap[row.EventID] = append(priceMap[row.EventID], row)
	}

	return priceMap, nil
}
