package services

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/types"
)

type EventRegistrationService struct {
	db           queries.DBTX
	eventService *EventService
}

func NewEventRegistrationService(db queries.DBTX) *EventRegistrationService {
	return &EventRegistrationService{
		db:           db,
		eventService: NewEventService(db),
	}
}

func (s *EventRegistrationService) CreateEventRegistration(ctx context.Context, user *queries.User, params types.CreateEventRegistrationParams) (*queries.EventRegistration, error) {
	_, err := s.eventService.GetEventById(ctx, params.EventID)
	if err != nil {
		return nil, fmt.Errorf("CreateEventRegistration: %w", err)
	}

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("CreateEventRegistration: validation failed: %w", err)
	}

	tx, err := s.db.(*pgxpool.Pool).BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("CreateEventRegistration: %w", err)
	}

	return nil, nil
}
