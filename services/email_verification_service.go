package services

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/internal/crypto"
	"github.com/moroz/homeosapiens-go/internal/jobs"
)

type EmailVerificationService interface {
	VerifyEmailAddress(ctx context.Context, token string) (*queries.User, error)
	MaybeResendVerificationEmail(ctx context.Context, email string) (sent bool, err error)
}

type emailVerificationService struct {
	db queries.DBTX
}

func NewEmailVerificationService(db queries.DBTX) EmailVerificationService {
	return &emailVerificationService{db}
}

func (s *emailVerificationService) VerifyEmailAddress(ctx context.Context, token string) (*queries.User, error) {
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

// MaybeResendVerificationEmail checks the rate limit then enqueues a resend job.
// Returns false without an error when the email is already verified or not found, to avoid enumeration.
func (s *emailVerificationService) MaybeResendVerificationEmail(ctx context.Context, email string) (bool, error) {
	db := s.db.(*pgxpool.Pool)
	emailHash := crypto.HashEmail(email)

	result, err := queries.New(db).CheckUserEmailVerificationRateLimit(ctx, &queries.CheckUserEmailVerificationRateLimitParams{
		EmailHash: emailHash,
		RateLimitPeriod: pgtype.Interval{
			Microseconds: int64(config.EmailVerificationRateLimitPeriod / time.Microsecond),
			Valid:        true,
		},
	})
	if err != nil {
		return false, err
	}
	if !result.CanRequest {
		return false, errRateLimited{limitedUntil: result.LimitedUntil}
	}

	user, err := queries.New(db).GetUnverifiedUserByEmail(ctx, emailHash)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	tx, err := db.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer tx.Rollback(ctx)

	river, err := jobs.NewClient(db)
	if err != nil {
		return false, err
	}

	if _, err = river.InsertTx(ctx, tx, &jobs.SendUserEmailArgs{UserID: user.ID}, nil); err != nil {
		return false, err
	}

	if err := tx.Commit(ctx); err != nil {
		return false, err
	}

	return true, nil
}
