package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/internal/crypto"
	"github.com/moroz/homeosapiens-go/internal/jobs"
	"github.com/moroz/homeosapiens-go/types"
)

type UserPasswordResetService struct {
	db queries.DBTX
}

func NewUserPasswordResetService(db queries.DBTX) *UserPasswordResetService {
	return &UserPasswordResetService{db}
}

var ErrUserNonExistent = fmt.Errorf("user not found")

func (s *UserPasswordResetService) MaybeIssuePasswordResetTokenForUser(ctx context.Context, email string) (bool, error) {
	emailHash := crypto.HashEmail(email)

	user, err := queries.New(s.db).GetUserByEmail(ctx, emailHash)
	if errors.Is(err, sql.ErrNoRows) {
		return false, ErrUserNonExistent
	}
	if err != nil {
		return false, nil
	}

	result, err := queries.New(s.db).CheckUserTokenFlowRateLimit(ctx, &queries.CheckUserTokenFlowRateLimitParams{
		EmailHash: emailHash,
		Context:   config.UserTokenContextPasswordReset,
		RateLimitPeriod: pgtype.Interval{
			Microseconds: int64(config.PasswordResetRateLimitPeriod / time.Microsecond),
			Valid:        true,
		},
	})
	if err != nil {
		return false, err
	}
	if !result.CanRequest {
		return false, errRateLimited{limitedUntil: result.LimitedUntil}
	}

	river, err := jobs.NewClient(s.db)
	if err != nil {
		return false, err
	}

	_, err = river.Insert(ctx, &jobs.SendUserEmailArgs{
		UserID:    user.ID,
		EmailType: jobs.UserEmailTypePasswordReset,
	}, nil)

	return err == nil, err
}

func (s *UserPasswordResetService) IssuePasswordResetTokenForUser(ctx context.Context, user *queries.User) (*types.UserTokenDTO, error) {
	return NewUserTokenService(s.db).IssueHashedTokenForUser(ctx, user, config.UserTokenContextPasswordReset, config.PasswordResetTokenValidity)
}
