package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/internal/jobs"
)

type EventRegistrationService struct {
	db *pgxpool.Pool
}

func NewEventRegistrationService(db *pgxpool.Pool) *EventRegistrationService {
	return &EventRegistrationService{db: db}
}

func (s *EventRegistrationService) CreateEventRegistration(ctx context.Context, user *queries.User, event *queries.Event) (bool, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return false, fmt.Errorf("CreateEventRegistration: %w", err)
	}
	defer tx.Rollback(ctx)

	result, err := queries.New(tx).InsertEventRegistration(ctx, &queries.InsertEventRegistrationParams{
		EventID: event.ID,
		UserID:  user.ID,
	})

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, err
	}

	newRegistration := result != nil

	if newRegistration {
		river, err := jobs.NewClient(s.db)
		if err != nil {
			return false, fmt.Errorf("CreateEventRegistration: %w", err)
		}
		_, err = river.InsertTx(ctx, tx, &jobs.SendEventRegistrationEmailArgs{
			UserID:  user.ID,
			EventID: event.ID,
		}, nil)
		if err != nil {
			return false, fmt.Errorf("CreateEventRegistration: failed to enqueue confirmation email: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return false, fmt.Errorf("CreateEventRegistration: %w", err)
	}

	return newRegistration, nil
}

func (s *EventRegistrationService) DeleteEventRegistration(ctx context.Context, user *queries.User, event *queries.Event) (bool, error) {
	_, err := queries.New(s.db).DeleteEventRegistration(ctx, &queries.DeleteEventRegistrationParams{
		EventID: event.ID,
		UserID:  user.ID,
	})

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}
