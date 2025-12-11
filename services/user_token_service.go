package services

import (
	"context"
	"crypto/rand"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/moroz/homeosapiens-go/db/queries"
)

const TokenLength = 32
const UserTokenDefaultValidity = 24 * time.Hour

type UserTokenService struct {
	db queries.DBTX
}

func NewUserTokenService(db queries.DBTX) *UserTokenService {
	return &UserTokenService{db}
}

func (s *UserTokenService) IssueAccessTokenForUser(ctx context.Context, user *queries.User, validity time.Duration) (*queries.UserToken, error) {
	var token = make([]byte, TokenLength)
	if _, err := rand.Read(token); err != nil {
		return nil, err
	}

	return queries.New(s.db).InsertUserToken(ctx, &queries.InsertUserTokenParams{
		UserID:  user.ID,
		Context: "access",
		Token:   token,
		ValidUntil: pgtype.Timestamp{
			Valid: true,
			Time:  time.Now().Add(validity),
		},
	})
}

func (s *UserTokenService) DeleteToken(ctx context.Context, token []byte) (bool, error) {
	return queries.New(s.db).DeleteUserToken(ctx, token)
}
