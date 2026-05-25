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
	"github.com/moroz/homeosapiens-go/types"
)

type EmailVerificationService interface {
	VerifyEmailAddress(ctx context.Context, token string) (*queries.User, error)
	IssueEmailVerificationTokenForUser(ctx context.Context, user *queries.User) (*types.UserTokenDTO, error)
	MaybeResendVerificationEmail(ctx context.Context, email string) (sent bool, err error)
}

type emailVerificationService struct {
	db *pgxpool.Pool
}

func NewEmailVerificationService(db *pgxpool.Pool) EmailVerificationService {
	return &emailVerificationService{db}
}

func (s *emailVerificationService) VerifyEmailAddress(ctx context.Context, token string) (*queries.User, error) {
	binary, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return nil, fmt.Errorf("VerifyEmailAddress: %w", err)
	}

	tx, err := s.db.Begin(ctx)
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

// IssueEmailVerificationTokenForUser is called from the background worker.
func (s *emailVerificationService) IssueEmailVerificationTokenForUser(ctx context.Context, user *queries.User) (*types.UserTokenDTO, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	if user.EmailConfirmedAt != nil {
		return nil, ErrEmailAlreadyVerified
	}

	if err := queries.New(s.db).DeletePreexistingEmailVerificationTokens(ctx, user.ID); err != nil {
		return nil, err
	}

	token, err := NewUserTokenService(s.db).IssueHashedTokenForUser(ctx, user, "email_verification", config.EmailVerificationTokenValidity)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return token, nil
}

// MaybeResendVerificationEmail checks the rate limit then enqueues a resend job.
// Returns false without an error when the email is already verified or not found, to avoid enumeration.
func (s *emailVerificationService) MaybeResendVerificationEmail(ctx context.Context, email string) (bool, error) {
	emailHash := crypto.HashEmail(email)

	result, err := queries.New(s.db).CheckUserTokenFlowRateLimit(ctx, &queries.CheckUserTokenFlowRateLimitParams{
		EmailHash: emailHash,
		Context:   "email_verification",
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

	user, err := queries.New(s.db).GetUnverifiedUserByEmail(ctx, emailHash)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	river, err := jobs.NewClient(s.db)
	if err != nil {
		return false, err
	}

	_, err = river.Insert(ctx, &jobs.SendUserEmailArgs{
		UserID:    user.ID,
		EmailType: jobs.UserEmailTypeEmailVerification,
	}, nil)

	return err == nil, err
}
