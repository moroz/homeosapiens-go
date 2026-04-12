package services

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/alexedwards/argon2id"
	"github.com/bincyber/go-sqlcrypter"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/internal/crypto"
	"github.com/moroz/homeosapiens-go/internal/jobs"
	"github.com/moroz/homeosapiens-go/types"
)

type UserRegistrationService struct {
	db queries.DBTX
}

func NewUserRegistrationService(db queries.DBTX) *UserRegistrationService {
	return &UserRegistrationService{
		db: db,
	}
}

func (s *UserRegistrationService) RegisterUser(ctx context.Context, params *types.RegisterUserParams) (*queries.User, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	passwordHash, err := argon2id.CreateHash(params.Password, config.ResolveArgon2Params())
	if err != nil {
		return nil, err
	}

	tx, err := s.db.(*pgxpool.Pool).Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("RegisterUser: %w", err)
	}
	defer tx.Rollback(ctx)

	params.Email = strings.TrimSpace(strings.ToLower(params.Email))

	user, err := queries.New(tx).InsertUser(ctx, &queries.InsertUserParams{
		PreferredLocale: queries.Locale(params.PreferredLocale),
		Email:           sqlcrypter.NewEncryptedBytes(params.Email),
		EmailHash:       crypto.HashEmail(params.Email),
		GivenName:       sqlcrypter.NewEncryptedBytes(params.GivenName),
		FamilyName:      sqlcrypter.NewEncryptedBytes(params.FamilyName),
		PasswordHash:    &passwordHash,
	})

	if err, ok := errors.AsType[*pgconn.PgError](err); ok && err.Code == "23505" && err.ConstraintName == "users_email_hash_key" {
		return nil, validation.Errors{
			"email": validation.NewError("unique", "has already been taken"),
		}
	}
	if err != nil {
		return nil, fmt.Errorf("RegisterUser: %w", err)
	}

	river, err := jobs.NewClient(s.db)
	if err != nil {
		return nil, fmt.Errorf("RegisterUser: %w", err)
	}

	_, err = river.InsertTx(ctx, tx, &jobs.SendUserEmailArgs{
		UserID: user.ID,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("RegisterUser: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("RegisterUser: %w", err)
	}

	return user, nil
}

var ErrTokenInvalid = errors.New("VerifyEmailAddress: token invalid or expired")

func (s *UserRegistrationService) VerifyEmailAddress(ctx context.Context, token string) (*queries.User, error) {
	binary, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return nil, fmt.Errorf("VerifyEmailAddress: %w", err)
	}

	tx, err := s.db.(*pgxpool.Pool).Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("VerifyEmailAddress: %w", err)
	}
	defer tx.Rollback(ctx)

	hashed := crypto.HashUserToken(binary)

	user, err := queries.New(tx).VerifyEmailAddressByUserToken(ctx, hashed)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrTokenInvalid
	}
	if err != nil {
		return nil, fmt.Errorf("VerifyEmailAddress: %w", err)
	}

	if _, err := queries.New(tx).DeleteUserToken(ctx, hashed); err != nil {
		return nil, fmt.Errorf("VerifyEmailAddress: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("VerifyEmailAddress: %w", err)
	}

	return user, nil
}
