package services

import (
	"context"
	"fmt"

	"github.com/bincyber/go-sqlcrypter"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/types"
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

func (s *EventRegistrationService) CreateEventRegistration(ctx context.Context, user *queries.User, params *types.CreateEventRegistrationParams) (*queries.EventRegistration, error) {
	event, err := s.eventService.GetEventById(ctx, params.EventID)
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

	if user == nil {
		user, err = s.userService.FindOrCreateUserFromEventRegistrationParams(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("CreateEventRegistration: %w", err)
		}
	}

	registration, err := queries.New(s.db).InsertEventRegistration(ctx, &queries.InsertEventRegistrationParams{
		EventID:    event.ID,
		UserID:     user.ID,
		GivenName:  sqlcrypter.NewEncryptedBytes(params.GivenName),
		FamilyName: sqlcrypter.NewEncryptedBytes(params.FamilyName),
		Email:      sqlcrypter.NewEncryptedBytes(params.Email),
		Country:    params.Country,
	})

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("CreateEventRegistration: %w", err)
	}

	return registration, nil
}
