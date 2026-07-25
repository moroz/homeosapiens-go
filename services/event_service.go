package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/shopspring/decimal"
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

func (s *EventService) GetPaidEventById(ctx context.Context, id uuid.UUID) (*queries.GetPaidEventByIdRow, error) {
	return queries.New(s.db).GetPaidEventById(ctx, id)
}

type EventListDto struct {
	*queries.ListEventsRow
	Product           *queries.Product
	Hosts             []*queries.ListHostsForEventsRow
	Prices            []*queries.ProductPrice
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

	products, err := s.preloadProductsForEvents(ctx, ids)
	if err != nil {
		return nil, err
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
			Product:           products[event.ID],
			Hosts:             hosts[event.ID],
			Prices:            prices[event.ID],
			EventRegistration: registrations[event.ID],
			RegistrationCount: regCounts[event.ID],
			CountInCart:       cartCounts[event.ID],
		})
	}

	return result, nil
}

func (s *EventService) GetEventDetailsById(ctx context.Context, eventId uuid.UUID, user *queries.User) (*types.EventDetailsDto, error) {
	event, err := queries.New(s.db).GetEventById(ctx, eventId)
	if err != nil {
		return nil, err
	}

	return s.GetEventDetailsForEvent(ctx, event, user, nil)
}

func (s *EventService) GetEventDetailsBySlug(ctx context.Context, slug string, user *queries.User, cartId *uuid.UUID) (*types.EventDetailsDto, error) {
	event, err := queries.New(s.db).GetEventBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	return s.GetEventDetailsForEvent(ctx, event, user, cartId)
}

func (s *EventService) GetEventDetailsForEvent(ctx context.Context, event *queries.Event, user *queries.User, cartId *uuid.UUID) (*types.EventDetailsDto, error) {
	var dto types.EventDetailsDto
	dto.Event = event

	products, err := s.preloadProductsForEvents(ctx, []uuid.UUID{event.ID})
	if err != nil {
		return nil, err
	}
	dto.Product = products[event.ID]

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

func (s *EventService) preloadProductsForEvents(ctx context.Context, eventIds []uuid.UUID) (map[uuid.UUID]*queries.Product, error) {
	products, err := queries.New(s.db).ListProductsForEvents(ctx, eventIds)
	if err != nil {
		return nil, err
	}

	productMap := make(map[uuid.UUID]*queries.Product)
	for _, row := range products {
		productMap[row.EventID] = &row.Product
	}

	return productMap, nil
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

func (s *EventService) preloadPricesForEvents(ctx context.Context, eventIds []uuid.UUID) (map[uuid.UUID][]*queries.ProductPrice, error) {
	prices, err := queries.New(s.db).ListPricesForEvents(ctx, eventIds)
	if err != nil {
		return nil, err
	}

	priceMap := make(map[uuid.UUID][]*queries.ProductPrice)
	for _, row := range prices {
		priceMap[row.EventID] = append(priceMap[row.EventID], &row.ProductPrice)
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

func (s *EventService) preloadCartLineItemPresenceForEvents(ctx context.Context, cartID *uuid.UUID, eventIDs []uuid.UUID) (map[uuid.UUID]int, error) {
	result := make(map[uuid.UUID]int)
	if cartID == nil {
		return result, nil
	}

	counts, err := queries.New(s.db).CountCartLineItemQuantitiesForProducts(ctx, &queries.CountCartLineItemQuantitiesForProductsParams{
		EventIds: eventIDs,
		CartID:   *cartID,
	})
	if err != nil {
		return result, nil
	}

	for _, row := range counts {
		result[row.EventID] = int(row.Quantity)
	}
	return result, nil
}

func (s *EventService) CreateEvent(ctx context.Context, params *types.CreateEventInput) (*types.EventDetailsDto, error) {
	tx, err := s.db.(*pgxpool.Pool).Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var productId *uuid.UUID
	var product *queries.Product

	if params.Price != nil && !params.Price.Equal(decimal.Zero) {
		product, err = queries.New(tx).InsertProduct(ctx, &queries.InsertProductParams{
			ProductType:       queries.ProductTypeEvent,
			TitlePl:           params.TitlePl,
			TitleEn:           params.TitleEn,
			BasePriceAmount:   *params.Price,
			BasePriceCurrency: *params.Currency,
		})
		if err != nil {
			return nil, err
		}
		productId = &product.ID
	}

	event, err := queries.New(tx).InsertEvent(ctx, &queries.InsertEventParams{
		TitleEn:          params.TitleEn,
		TitlePl:          params.TitlePl,
		StartsAt:         params.StartsAt,
		EndsAt:           params.EndsAt,
		IsVirtual:        params.IsVirtual,
		DescriptionEn:    params.DescriptionEn,
		DescriptionPl:    params.DescriptionPl,
		EventType:        queries.EventType(params.EventType),
		Slug:             params.Slug,
		SubtitleEn:       params.SubtitleEn,
		SubtitlePl:       params.SubtitlePl,
		VenueNameEn:      params.VenueNameEn,
		VenueNamePl:      params.VenueNamePl,
		VenueStreet:      params.VenueStreet,
		VenueCityEn:      params.VenueCityEn,
		VenueCityPl:      params.VenueCityPl,
		VenuePostalCode:  params.VenuePostalCode,
		VenueCountryCode: params.VenueCountryCode,
		ProductID:        productId,
	})
	if err != nil {
		return nil, err
	}

	for i, hostId := range params.HostIds {
		_, err := queries.New(tx).InsertEventHost(ctx, &queries.InsertEventHostParams{
			EventID:  event.ID,
			HostID:   hostId,
			Position: int32(i + 1),
		})
		if err != nil {
			return nil, err
		}
	}

	hosts, err := s.preloadHostsForEvents(ctx, []uuid.UUID{event.ID})
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &types.EventDetailsDto{
		Event:   event,
		Product: product,
		Prices:  nil,
		Hosts:   hosts[event.ID],
	}, nil
}

func (s *EventService) UpdateEvent(ctx context.Context, eventId uuid.UUID, params *types.UpdateEventInput) (*queries.Event, error) {
	tx, err := s.db.(*pgxpool.Pool).Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	return queries.New(tx).UpdateEvent(ctx, &queries.UpdateEventParams{
		TitlePl:       params.TitlePl,
		TitleEn:       params.TitleEn,
		SubtitlePl:    params.SubtitlePl,
		SubtitleEn:    params.SubtitleEn,
		DescriptionPl: params.DescriptionPl,
		DescriptionEn: params.DescriptionEn,
		EventID:       eventId,
	})
}
