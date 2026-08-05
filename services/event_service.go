package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

type EventListDto struct {
	*queries.ListPublishedEventsRow
	Product           *queries.Product
	Hosts             []*queries.ListHostsForEventsRow
	Prices            []*queries.ProductPrice
	EventRegistration *queries.EventRegistration
	RegistrationCount int
	CountInCart       int
}

func (s *EventService) ListPublishedEventsForUser(ctx context.Context, user *queries.User, cartId *uuid.UUID) ([]*EventListDto, error) {
	events, err := queries.New(s.db).ListPublishedEvents(ctx)
	if err != nil {
		return nil, err
	}

	var ids []uuid.UUID
	for _, event := range events {
		ids = append(ids, event.ID)
	}

	products, err := s.preloadProductsForEvents(ctx, s.db, ids)
	if err != nil {
		return nil, err
	}

	hosts, err := s.preloadHostsForEvents(ctx, s.db, ids)
	if err != nil {
		return nil, err
	}

	prices, err := s.preloadPricesForEvents(ctx, s.db, ids)
	if err != nil {
		return nil, err
	}

	registrations, err := s.preloadEventRegistrationsForEvents(ctx, s.db, ids, user)
	if err != nil {
		return nil, err
	}

	regCounts, err := s.preloadRegistrationCountsForEvents(ctx, s.db, ids)
	if err != nil {
		return nil, err
	}

	cartCounts, err := s.preloadCartLineItemPresenceForEvents(ctx, s.db, cartId, ids)
	if err != nil {
		return nil, err
	}

	var result []*EventListDto
	for _, event := range events {
		result = append(result, &EventListDto{
			ListPublishedEventsRow: event,
			Product:                products[event.ID],
			Hosts:                  hosts[event.ID],
			Prices:                 prices[event.ID],
			EventRegistration:      registrations[event.ID],
			RegistrationCount:      regCounts[event.ID],
			CountInCart:            cartCounts[event.ID],
		})
	}

	return result, nil
}

func (s *EventService) GetEventDetailsById(ctx context.Context, eventId uuid.UUID, user *queries.User) (*types.EventDetailsDto, error) {
	event, err := queries.New(s.db).GetEventById(ctx, eventId)
	if err != nil {
		return nil, err
	}

	return s.eventDetailsForEvent(ctx, s.db, event, user, nil)
}

func (s *EventService) GetEventDetailsBySlug(ctx context.Context, slug string, user *queries.User, cartId *uuid.UUID) (*types.EventDetailsDto, error) {
	event, err := queries.New(s.db).GetEventBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	return s.eventDetailsForEvent(ctx, s.db, event, user, cartId)
}

// GetEventDetailsForEvent loads an event's associations off the pool.
func (s *EventService) GetEventDetailsForEvent(ctx context.Context, event *queries.Event, user *queries.User, cartId *uuid.UUID) (*types.EventDetailsDto, error) {
	return s.eventDetailsForEvent(ctx, s.db, event, user, cartId)
}

// eventDetailsForEvent reads through db, so a caller holding a transaction can
// build the DTO from its own uncommitted writes.
func (s *EventService) eventDetailsForEvent(ctx context.Context, db queries.DBTX, event *queries.Event, user *queries.User, cartId *uuid.UUID) (*types.EventDetailsDto, error) {
	var dto types.EventDetailsDto
	dto.Event = event

	products, err := s.preloadProductsForEvents(ctx, db, []uuid.UUID{event.ID})
	if err != nil {
		return nil, err
	}
	dto.Product = products[event.ID]

	prices, err := s.preloadPricesForEvents(ctx, db, []uuid.UUID{event.ID})
	if err != nil {
		return nil, err
	}
	dto.Prices = prices[event.ID]

	hosts, err := s.preloadHostsForEvents(ctx, db, []uuid.UUID{event.ID})
	if err != nil {
		return nil, err
	}
	dto.Hosts = hosts[event.ID]

	registrations, err := s.preloadEventRegistrationsForEvents(ctx, db, []uuid.UUID{event.ID}, user)
	if err != nil {
		return nil, err
	}
	dto.EventRegistration = registrations[event.ID]

	counts, err := s.preloadRegistrationCountsForEvents(ctx, db, []uuid.UUID{event.ID})
	if err != nil {
		return nil, err
	}
	dto.RegistrationCount = counts[event.ID]

	if cartId != nil {
		cartCounts, err := s.preloadCartLineItemPresenceForEvents(ctx, db, cartId, []uuid.UUID{event.ID})
		if err != nil {
			return nil, err
		}
		dto.CountInCart = cartCounts[event.ID]
	}

	return &dto, nil
}

func (s *EventService) preloadProductsForEvents(ctx context.Context, db queries.DBTX, eventIds []uuid.UUID) (map[uuid.UUID]*queries.Product, error) {
	products, err := queries.New(db).ListProductsForEvents(ctx, eventIds)
	if err != nil {
		return nil, err
	}

	productMap := make(map[uuid.UUID]*queries.Product)
	for _, row := range products {
		productMap[row.EventID] = &row.Product
	}

	return productMap, nil
}

func (s *EventService) preloadHostsForEvents(ctx context.Context, db queries.DBTX, eventIds []uuid.UUID) (map[uuid.UUID][]*queries.ListHostsForEventsRow, error) {
	hosts, err := queries.New(db).ListHostsForEvents(ctx, eventIds)
	if err != nil {
		return nil, err
	}

	hostMap := make(map[uuid.UUID][]*queries.ListHostsForEventsRow)
	for _, row := range hosts {
		hostMap[row.EventID] = append(hostMap[row.EventID], row)
	}

	return hostMap, nil
}

func (s *EventService) preloadPricesForEvents(ctx context.Context, db queries.DBTX, eventIds []uuid.UUID) (map[uuid.UUID][]*queries.ProductPrice, error) {
	prices, err := queries.New(db).ListPricesForEvents(ctx, eventIds)
	if err != nil {
		return nil, err
	}

	priceMap := make(map[uuid.UUID][]*queries.ProductPrice)
	for _, row := range prices {
		priceMap[row.EventID] = append(priceMap[row.EventID], &row.ProductPrice)
	}

	return priceMap, nil
}

func (s *EventService) preloadEventRegistrationsForEvents(ctx context.Context, db queries.DBTX, eventIds []uuid.UUID, user *queries.User) (map[uuid.UUID]*queries.EventRegistration, error) {
	resultMap := make(map[uuid.UUID]*queries.EventRegistration)

	if user == nil {
		return resultMap, nil
	}

	registrations, err := queries.New(db).ListEventRegistrationsForUserForEvents(ctx, &queries.ListEventRegistrationsForUserForEventsParams{
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

func (s *EventService) preloadRegistrationCountsForEvents(ctx context.Context, db queries.DBTX, eventIds []uuid.UUID) (map[uuid.UUID]int, error) {
	counts, err := queries.New(db).CountRegistrationsForEvents(ctx, eventIds)
	if err != nil {
		return nil, err
	}

	result := make(map[uuid.UUID]int)
	for _, row := range counts {
		result[row.EventID] = int(row.Count)
	}
	return result, nil
}

func (s *EventService) preloadCartLineItemPresenceForEvents(ctx context.Context, db queries.DBTX, cartID *uuid.UUID, eventIDs []uuid.UUID) (map[uuid.UUID]int, error) {
	result := make(map[uuid.UUID]int)
	if cartID == nil {
		return result, nil
	}

	counts, err := queries.New(db).CountCartLineItemQuantitiesForProducts(ctx, &queries.CountCartLineItemQuantitiesForProductsParams{
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
	if err := params.Validate(); err != nil {
		return nil, err
	}

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
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok && pgErr.Code == "23505" && pgErr.ConstraintName == "events_slug_idx" {
		return nil, validation.Errors{
			"slug": validation.NewError("unique", "has already been taken"),
		}
	}
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

	// Read the hosts back through tx: the inserts above are not visible on the
	// pool until the transaction commits.
	hosts, err := s.preloadHostsForEvents(ctx, tx, []uuid.UUID{event.ID})
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

func (s *EventService) PublishEvent(ctx context.Context, eventId uuid.UUID) (*queries.Event, error) {
	event, err := queries.New(s.db).GetEventById(ctx, eventId)
	if err != nil {
		return nil, err
	}

	params := &types.PublishEventValidation{
		DescriptionPl: event.DescriptionPl,
		DescriptionEn: event.DescriptionEn,
		PublishedAt:   event.PublishedAt,
	}
	if err := params.Validate(); err != nil {
		return nil, err
	}

	return queries.New(s.db).PublishEvent(ctx, eventId)
}

// UnpublishEvent takes an event off the public listing. Unlike deletion it is
// always allowed: existing registrations keep their record of the event, they
// just stop being joined by new sign-ups. Unpublishing an event that is already
// a draft is a no-op.
func (s *EventService) UnpublishEvent(ctx context.Context, eventId uuid.UUID) (*queries.Event, error) {
	return queries.New(s.db).UnpublishEvent(ctx, eventId)
}

// DeleteEvent removes an event along with its host associations. It refuses to
// delete an event anyone has signed up for, returning ErrEventHasRegistrations;
// unpublishing is the way to retire such an event. The event's product, if any,
// is left behind, because order line items still reference it.
func (s *EventService) DeleteEvent(ctx context.Context, eventId uuid.UUID) error {
	tx, err := s.db.(*pgxpool.Pool).Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Locking the row keeps a registration from landing between the count and
	// the delete; the registration path reads the event before inserting.
	var locked uuid.UUID
	err = tx.QueryRow(ctx, "select id from events where id = $1 for update", eventId).Scan(&locked)
	if err != nil {
		return err
	}

	count, err := queries.New(tx).CountRegistrationsForEvent(ctx, eventId)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrEventHasRegistrations
	}

	// events_hosts has no cascade, so the join rows go first.
	if err := queries.New(tx).DeleteEventHosts(ctx, eventId); err != nil {
		return err
	}

	deleted, err := queries.New(tx).DeleteEvent(ctx, eventId)
	if err != nil {
		return err
	}
	if deleted == 0 {
		return sql.ErrNoRows
	}

	return tx.Commit(ctx)
}

// UpdateEvent applies a selective update to an event. This bespoke logic is
// intended to mimic the selective update behavior of Ecto.Changeset in Elixir:
// only the fields present in the payload are written, and fields backed by a
// nullable column can be cleared by passing an explicit null.
func (s *EventService) UpdateEvent(ctx context.Context, eventId uuid.UUID, params *types.PatchEventInput) (*types.EventDetailsDto, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	event, err := s.GetEventById(ctx, eventId)
	if err != nil {
		return nil, err
	}

	// The venue requirement depends on the row being patched, so it can only be
	// checked once the event is loaded.
	if err := params.ValidateVenue(event); err != nil {
		return nil, err
	}

	if params.IsEmpty() {
		return s.eventDetailsForEvent(ctx, s.db, event, nil, nil)
	}

	tx, err := s.db.(*pgxpool.Pool).Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// Pricing may attach a brand-new product to the event, in which case the
	// events row has to carry the resulting foreign key.
	productId, err := s.patchEventProduct(ctx, tx, event, params)
	if err != nil {
		return nil, err
	}

	// The column list is the contract between the payload and the table, so it
	// is spelled out rather than derived from field names at runtime.
	assignments := []struct {
		column string
		value  types.Optional[string]
	}{
		{"title_en", params.TitleEn},
		{"title_pl", params.TitlePl},
		{"slug", params.Slug},
		{"subtitle_en", params.SubtitleEn},
		{"subtitle_pl", params.SubtitlePl},
		{"description_en", params.DescriptionEn},
		{"description_pl", params.DescriptionPl},
		{"venue_name_en", params.VenueNameEn},
		{"venue_name_pl", params.VenueNamePl},
		{"venue_street", params.VenueStreet},
		{"venue_city_en", params.VenueCityEn},
		{"venue_city_pl", params.VenueCityPl},
		{"venue_postal_code", params.VenuePostalCode},
		{"venue_country_code", params.VenueCountryCode},
	}

	timestampAssignments := []struct {
		column string
		value  types.Optional[time.Time]
	}{
		{"starts_at", params.StartsAt},
		{"ends_at", params.EndsAt},
	}

	boolAssignments := []struct {
		column string
		value  types.Optional[bool]
	}{
		{"is_virtual", params.IsVirtual},
	}

	var query strings.Builder
	query.WriteString("update events set ")

	var queryVars []any
	for _, assignment := range assignments {
		if !assignment.value.Set {
			continue
		}

		queryVars = append(queryVars, assignment.value.Ptr())
		fmt.Fprintf(&query, "%s = $%d, ", assignment.column, len(queryVars))
	}

	for _, assignment := range timestampAssignments {
		if !assignment.value.Set {
			continue
		}

		queryVars = append(queryVars, assignment.value.Ptr())
		fmt.Fprintf(&query, "%s = $%d, ", assignment.column, len(queryVars))
	}

	for _, assignment := range boolAssignments {
		if !assignment.value.Set {
			continue
		}

		queryVars = append(queryVars, assignment.value.Ptr())
		fmt.Fprintf(&query, "%s = $%d, ", assignment.column, len(queryVars))
	}

	if productId != nil {
		queryVars = append(queryVars, productId)
		fmt.Fprintf(&query, "product_id = $%d, ", len(queryVars))
	}

	queryVars = append(queryVars, eventId)
	fmt.Fprintf(&query, "updated_at = now() where id = $%d", len(queryVars))

	if _, err := tx.Exec(ctx, query.String(), queryVars...); err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok && pgErr.Code == "23505" && pgErr.ConstraintName == "events_slug_idx" {
			return nil, validation.Errors{
				"slug": validation.NewError("unique", "has already been taken"),
			}
		}
		return nil, err
	}

	if params.HostIds.Set {
		if err := s.replaceEventHosts(ctx, tx, eventId, params.HostIds.Value); err != nil {
			return nil, err
		}
	}

	updated, err := queries.New(tx).GetEventById(ctx, eventId)
	if err != nil {
		return nil, err
	}

	// Built from tx, so the caller sees the associations this update just wrote
	// without a second round trip outside the transaction.
	details, err := s.eventDetailsForEvent(ctx, tx, updated, nil, nil)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return details, nil
}

// patchEventProduct applies a pricing change to the event's product, returning
// the product ID that the events row must be pointed at, or nil when the
// existing association already holds.
func (s *EventService) patchEventProduct(ctx context.Context, tx pgx.Tx, event *queries.Event, params *types.PatchEventInput) (*uuid.UUID, error) {
	if !params.Price.Set && !params.Currency.Set {
		return nil, nil
	}

	// An event that already has a product keeps it, even when the price drops
	// to zero: cart and order line items reference it by ID.
	if event.ProductID != nil {
		product, err := queries.New(tx).GetProductById(ctx, *event.ProductID)
		if err != nil {
			return nil, err
		}

		amount := product.BasePriceAmount
		if params.Price.Set {
			amount = params.Price.Value
		}

		currency := product.BasePriceCurrency
		if params.Currency.Set {
			currency = params.Currency.Value
		}

		_, err = queries.New(tx).UpdateProductPrice(ctx, &queries.UpdateProductPriceParams{
			ProductID:         product.ID,
			BasePriceAmount:   amount,
			BasePriceCurrency: currency,
		})

		return nil, err
	}

	// A free event stays free until it is given a non-zero price; a currency on
	// its own has nowhere to be stored.
	if !params.Price.Set || params.Price.Value.Equal(decimal.Zero) {
		return nil, nil
	}

	if !params.Currency.Set {
		return nil, validation.Errors{
			"currency": validation.NewError("required", "is required when setting a price"),
		}
	}

	titlePl := event.TitlePl
	if params.TitlePl.Set {
		titlePl = params.TitlePl.Value
	}

	titleEn := event.TitleEn
	if params.TitleEn.Set {
		titleEn = params.TitleEn.Value
	}

	product, err := queries.New(tx).InsertProduct(ctx, &queries.InsertProductParams{
		ProductType:       queries.ProductTypeEvent,
		TitlePl:           titlePl,
		TitleEn:           titleEn,
		BasePriceAmount:   params.Price.Value,
		BasePriceCurrency: params.Currency.Value,
	})
	if err != nil {
		return nil, err
	}

	return &product.ID, nil
}

// replaceEventHosts swaps the event's entire host list. Positions are unique
// per event, so incremental edits would have to shuffle rows around a unique
// index; deleting and reinserting inside the caller's transaction is simpler.
func (s *EventService) replaceEventHosts(ctx context.Context, tx pgx.Tx, eventId uuid.UUID, hostIds []uuid.UUID) error {
	if err := queries.New(tx).DeleteEventHosts(ctx, eventId); err != nil {
		return err
	}

	for i, hostId := range hostIds {
		_, err := queries.New(tx).InsertEventHost(ctx, &queries.InsertEventHostParams{
			EventID:  eventId,
			HostID:   hostId,
			Position: int32(i + 1),
		})
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok && pgErr.Code == "23503" {
			return validation.Errors{
				"hostIds": validation.NewError("exists", "references a host that does not exist"),
			}
		}
		if err != nil {
			return err
		}
	}

	return nil
}
