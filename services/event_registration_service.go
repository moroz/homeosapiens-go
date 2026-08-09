package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jackc/pgx/v5"
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

func (s *EventRegistrationService) CreateEventRegistration(ctx context.Context, user *queries.User, event *queries.Event) (*queries.InsertEventRegistrationRow, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("CreateEventRegistration: %w", err)
	}
	defer tx.Rollback(ctx)

	queryResult, err := queries.New(tx).InsertEventRegistration(ctx, &queries.InsertEventRegistrationParams{
		EventID: event.ID,
		UserID:  user.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("CreateEventRegistration: %w", err)
	}

	if queryResult.NewRecord {
		if err := s.sendEventRegistrationEmail(ctx, tx, queryResult); err != nil {
			return nil, fmt.Errorf("CreateEventRegistration: failed to enqueue confirmation email: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("CreateEventRegistration: %w", err)
	}

	return queryResult, nil
}

func (s *EventRegistrationService) sendEventRegistrationEmail(ctx context.Context, tx pgx.Tx, registration *queries.InsertEventRegistrationRow) error {
	river, err := jobs.NewClient(s.db)
	if err != nil {
		return err
	}
	_, err = river.InsertTx(ctx, tx, &jobs.SendEventRegistrationEmailArgs{
		UserID:  registration.EventRegistration.UserID,
		EventID: registration.EventRegistration.EventID,
	}, nil)
	return err
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

func countPages(count int64, perPage int32) int32 {
	return int32(
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
			Page:       params.Page,
			PerPage:    params.PerPage,
			TotalPages: countPages(count, params.PerPage),
		},
		Data: rows,
	}, nil
}

func (s *EventRegistrationService) ListEligibleUsersForEvent(ctx context.Context, params *types.ListEligibleUsersForEventParams) ([]*queries.ListEligibleUsersForEventRow, error) {
	allUsers, err := queries.New(s.db).ListEligibleUsersForEvent(ctx, params.EventID)
	if err != nil {
		return nil, err
	}

	q := strings.TrimSpace(params.SearchTerm)

	if q == "" {
		if len(allUsers) < 20 {
			return allUsers, nil
		}

		return allUsers[0:20], nil
	}

	var filtered []*queries.ListEligibleUsersForEventRow
	for _, user := range allUsers {
		fullName := fmt.Sprintf("%s %s", user.GivenName.Plaintext(), user.FamilyName.Plaintext())
		if strings.Contains(strings.ToLower(fullName), q) || strings.Contains(user.Email.Plaintext(), q) {
			filtered = append(filtered, user)
		}
		if len(filtered) == 20 {
			break
		}
	}

	return filtered, nil
}

// AdminCreateEventRegistration allows an administrator to enroll any user for any event, paid or free. It differs from CreatEventRegistration in that CreateEventRegistration is called with a user fetched from the request context, and an administrator can simply specify a user by primary key.
func (s *EventRegistrationService) AdminCreateEventRegistration(ctx context.Context, params *types.EnrollStudentForEventInput) (*queries.InsertEventRegistrationRow, error) {
	event, err := queries.New(s.db).GetEventById(ctx, params.EventID)
	if err != nil {
		return nil, fmt.Errorf("AdminCreateEventRegistration: %w", err)
	}

	if event.PublishedAt == nil {
		return nil, validation.Errors{
			"eventId": validation.NewError("unpublished", "event must be published"),
		}
	}
	if event.EndsAt.Before(time.Now()) {
		return nil, validation.Errors{
			"eventId": validation.NewError("expired", "event has already ended"),
		}
	}

	user, err := queries.New(s.db).GetUserByID(ctx, params.UserID)
	if err != nil {
		return nil, fmt.Errorf("AdminCreateEventRegistration: %w", err)
	}

	return s.CreateEventRegistration(ctx, user, event)
}
