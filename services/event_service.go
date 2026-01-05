package services

import (
	"context"
	"fmt"

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
	return queries.New(s.db).GetEventById(ctx, id)
}

type EventListDto struct {
	*queries.ListEventsRow
	Hosts             []*queries.ListHostsForEventsRow
	Prices            []*queries.EventPrice
	EventRegistration *queries.EventRegistration
}

func (s *EventService) ListEvents(ctx context.Context, user *queries.User) ([]*EventListDto, error) {
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

	registrations, err := s.preloadEventRegistrationsForEvents(ctx, ids, user)
	if err != nil {
		return nil, err
	}

	var result []*EventListDto
	for _, event := range events {
		result = append(result, &EventListDto{
			ListEventsRow:     event,
			Hosts:             hosts[event.ID],
			Prices:            prices[event.ID],
			EventRegistration: registrations[event.ID],
		})
	}

	return result, nil
}

type EventDetailsDto struct {
	*queries.Event
	Venue             *queries.Venue
	Prices            []*queries.EventPrice
	Hosts             []*queries.ListHostsForEventsRow
	EventRegistration *queries.EventRegistration
}

func (s *EventService) GetEventDetailsById(ctx context.Context, eventId string, user *queries.User) (*EventDetailsDto, error) {
	event, err := queries.New(s.db).GetEventById(ctx, eventId)
	if err != nil {
		return nil, err
	}

	return s.GetEventDetailsForEvent(ctx, event, user)
}

func (s *EventService) GetEventDetailsBySlug(ctx context.Context, slug string, user *queries.User) (*EventDetailsDto, error) {
	event, err := queries.New(s.db).GetEventBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	return s.GetEventDetailsForEvent(ctx, event, user)
}

func (s *EventService) GetEventDetailsForEvent(ctx context.Context, event *queries.Event, user *queries.User) (*EventDetailsDto, error) {
	var dto EventDetailsDto
	dto.Event = event

	if event.VenueID.Valid {
		venue, err := queries.New(s.db).GetVenueById(ctx, event.VenueID)
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

	registrations, err := s.preloadEventRegistrationsForEvents(ctx, []pgtype.UUID{event.ID}, user)
	if err != nil {
		return nil, err
	}
	dto.EventRegistration = registrations[event.ID]

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

func (s *EventService) preloadEventRegistrationsForEvents(ctx context.Context, eventIds []pgtype.UUID, user *queries.User) (map[pgtype.UUID]*queries.EventRegistration, error) {
	resultMap := make(map[pgtype.UUID]*queries.EventRegistration)

	if user == nil {
		return resultMap, nil
	}

	registrations, err := queries.New(s.db).ListEventRegistrationsForUserForEvents(ctx, &queries.ListEventRegistrationsForUserForEventsParams{
		Eventids: eventIds,
		Userid:   user.ID,
	})
	if err != nil {
		return resultMap, fmt.Errorf("preloadEventRegistrationsForEvents: %w", err)
	}
	for _, row := range registrations {
		resultMap[row.EventID] = row
	}
	return resultMap, nil
}
