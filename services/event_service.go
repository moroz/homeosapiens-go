package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/moroz/homeosapiens-go/db/queries"
)

type EventService struct {
	db queries.DBTX
}

func NewEventService(db queries.DBTX) *EventService {
	return &EventService{db}
}

func (s *EventService) GetEventById(ctx context.Context, id uuid.UUID) (*queries.Event, error) {
	return queries.New(s.db).GetEventById(ctx, id)
}

func (s *EventService) GetRegisterableEventById(ctx context.Context, id uuid.UUID) (*queries.Event, error) {
	return queries.New(s.db).GetFreeEventById(ctx, id)
}

func (s *EventService) GetPaidEventById(ctx context.Context, id uuid.UUID) (*queries.Event, error) {
	return queries.New(s.db).GetPaidEventById(ctx, id)
}

type EventListDto struct {
	*queries.ListEventsRow
	Hosts             []*queries.ListHostsForEventsRow
	Prices            []*queries.EventPrice
	EventRegistration *queries.EventRegistration
	RegistrationCount int
	CountInCart       int
}

func (s *EventService) ListEvents(ctx context.Context, user *queries.User, cartId *uuid.UUID) ([]*EventListDto, error) {
	events, err := queries.New(s.db).ListEvents(ctx)
	if err != nil {
		return nil, err
	}

	var ids []uuid.UUID
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

	regCounts, err := s.preloadRegistrationCountsForEvents(ctx, ids)
	if err != nil {
		return nil, err
	}

	cartCounts, err := s.preloadCartLineItemPresenceForEvents(ctx, cartId, ids)
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
			RegistrationCount: regCounts[event.ID],
			CountInCart:       cartCounts[event.ID],
		})
	}

	return result, nil
}

type EventDetailsDto struct {
	*queries.Event
	Prices            []*queries.EventPrice
	Hosts             []*queries.ListHostsForEventsRow
	EventRegistration *queries.EventRegistration
	RegistrationCount int
	CountInCart       int
}

func (s *EventService) GetEventDetailsById(ctx context.Context, eventId uuid.UUID, user *queries.User) (*EventDetailsDto, error) {
	event, err := queries.New(s.db).GetEventById(ctx, eventId)
	if err != nil {
		return nil, err
	}

	return s.GetEventDetailsForEvent(ctx, event, user, nil)
}

func (s *EventService) GetEventDetailsBySlug(ctx context.Context, slug string, user *queries.User, cartId *uuid.UUID) (*EventDetailsDto, error) {
	event, err := queries.New(s.db).GetEventBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	return s.GetEventDetailsForEvent(ctx, event, user, cartId)
}

func (s *EventService) GetEventDetailsForEvent(ctx context.Context, event *queries.Event, user *queries.User, cartId *uuid.UUID) (*EventDetailsDto, error) {
	var dto EventDetailsDto
	dto.Event = event

	prices, err := s.preloadPricesForEvents(ctx, []uuid.UUID{event.ID})
	if err != nil {
		return nil, err
	}
	dto.Prices = prices[event.ID]

	hosts, err := s.preloadHostsForEvents(ctx, []uuid.UUID{event.ID})
	if err != nil {
		return nil, err
	}
	dto.Hosts = hosts[event.ID]

	registrations, err := s.preloadEventRegistrationsForEvents(ctx, []uuid.UUID{event.ID}, user)
	if err != nil {
		return nil, err
	}
	dto.EventRegistration = registrations[event.ID]

	counts, err := s.preloadRegistrationCountsForEvents(ctx, []uuid.UUID{event.ID})
	if err != nil {
		return nil, err
	}
	dto.RegistrationCount = counts[event.ID]

	if cartId != nil {
		cartCounts, err := s.preloadCartLineItemPresenceForEvents(ctx, cartId, []uuid.UUID{event.ID})
		if err != nil {
			return nil, err
		}
		dto.CountInCart = cartCounts[event.ID]
	}

	return &dto, nil

}

func (s *EventService) preloadHostsForEvents(ctx context.Context, eventIds []uuid.UUID) (map[uuid.UUID][]*queries.ListHostsForEventsRow, error) {
	hosts, err := queries.New(s.db).ListHostsForEvents(ctx, eventIds)
	if err != nil {
		return nil, err
	}

	hostMap := make(map[uuid.UUID][]*queries.ListHostsForEventsRow)
	for _, row := range hosts {
		hostMap[row.EventID] = append(hostMap[row.EventID], row)
	}

	return hostMap, nil
}

func (s *EventService) preloadPricesForEvents(ctx context.Context, eventIds []uuid.UUID) (map[uuid.UUID][]*queries.EventPrice, error) {
	prices, err := queries.New(s.db).ListPricesForEvents(ctx, eventIds)
	if err != nil {
		return nil, err
	}

	priceMap := make(map[uuid.UUID][]*queries.EventPrice)
	for _, row := range prices {
		priceMap[row.EventID] = append(priceMap[row.EventID], row)
	}

	return priceMap, nil
}

func (s *EventService) preloadEventRegistrationsForEvents(ctx context.Context, eventIds []uuid.UUID, user *queries.User) (map[uuid.UUID]*queries.EventRegistration, error) {
	resultMap := make(map[uuid.UUID]*queries.EventRegistration)

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

func (s *EventService) preloadRegistrationCountsForEvents(ctx context.Context, eventIds []uuid.UUID) (map[uuid.UUID]int, error) {
	counts, err := queries.New(s.db).CountRegistrationsForEvents(ctx, eventIds)
	if err != nil {
		return nil, err
	}

	result := make(map[uuid.UUID]int)
	for _, row := range counts {
		result[row.EventID] = int(row.Count)
	}
	return result, nil
}

func (s *EventService) preloadCartLineItemPresenceForEvents(ctx context.Context, cartId *uuid.UUID, eventIds []uuid.UUID) (map[uuid.UUID]int, error) {
	result := make(map[uuid.UUID]int)
	if cartId == nil {
		return result, nil
	}

	counts, err := queries.New(s.db).CountCartLineItemQuantitiesForEvents(ctx, &queries.CountCartLineItemQuantitiesForEventsParams{
		EventIds: eventIds,
		CartID:   *cartId,
	})
	if err != nil {
		return result, nil
	}

	for _, row := range counts {
		result[row.EventID] = int(row.Quantity)
	}
	return result, nil
}
