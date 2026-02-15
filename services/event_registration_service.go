package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/moroz/homeosapiens-go/db/queries"
)

type EventRegistrationService struct {
	db           queries.DBTX
	eventService *EventService
	userService  *UserService
}

func NewEventRegistrationService(db queries.DBTX) *EventRegistrationService {
	return &EventRegistrationService{
		db:           db,
		eventService: NewEventService(db),
		userService:  NewUserService(db),
	}
}

func (s *EventRegistrationService) CreateEventRegistration(ctx context.Context, user *queries.User, event *queries.Event) (bool, error) {
	result, err := queries.New(s.db).InsertEventRegistration(ctx, &queries.InsertEventRegistrationParams{
		EventID: event.ID,
		UserID:  user.ID,
	})

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, err
	}

	return result != nil, nil
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
