package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/internal/jobs"
	"github.com/moroz/homeosapiens-go/types"
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

func countPages(count int64, perPage int32) int {
	return int(
		math.Ceil(
			float64(count) / float64(perPage),
		),
	)
}

func (s *EventRegistrationService) PaginateEventAttendants(ctx context.Context, params *queries.PaginateEventRegistrationsParams) (*types.PaginationPage[*queries.PaginateEventRegistrationsRow], error) {
	_, err := queries.New(s.db).GetEventById(ctx, params.EventID)
	if err != nil {
		return nil, err
	}

	rows, err := queries.New(s.db).PaginateEventRegistrations(ctx, params)
	if err != nil {
		return nil, err
	}

	count, err := queries.New(s.db).CountEventRegistrations(ctx, params.EventID)
	if err != nil {
		return nil, err
	}

	return &types.PaginationPage[*queries.PaginateEventRegistrationsRow]{
		Pagination: types.Pagination{
			Page:       int(params.Page),
			PerPage:    int(params.PerPage),
			TotalPages: countPages(count, params.PerPage),
		},
		Data: rows,
	}, nil
}
