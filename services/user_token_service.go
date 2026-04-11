package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/internal/crypto"
	"github.com/moroz/homeosapiens-go/types"
)

const TokenLength = 32
const UserTokenDefaultValidity = 24 * time.Hour

type UserTokenService struct {
	db queries.DBTX
}

func NewUserTokenService(db queries.DBTX) *UserTokenService {
	return &UserTokenService{db}
}

func generateToken() ([]byte, error) {
	var token = make([]byte, TokenLength)
	_, err := rand.Read(token)
	return token, err
}

func (s *UserTokenService) IssueAccessTokenForUser(ctx context.Context, user *queries.User, validity time.Duration) (*queries.UserToken, error) {
	token, err := generateToken()
	if err != nil {
		return nil, err
	}

	return queries.New(s.db).InsertUserToken(ctx, &queries.InsertUserTokenParams{
		UserID:     user.ID,
		Context:    "access",
		Token:      token,
		ValidUntil: time.Now().Add(validity),
	})
}

func (s *UserTokenService) DeleteToken(ctx context.Context, token []byte) (bool, error) {
	return queries.New(s.db).DeleteUserToken(ctx, token)
}

func (s *UserTokenService) IssueEmailVerificationTokenForUser(ctx context.Context, user *queries.User, validity time.Duration) (*types.UserTokenDTO, error) {
	token, err := generateToken()
	if err != nil {
		return nil, fmt.Errorf("IssueEmailVerificationTokenForUser: %w", err)
	}

	userToken, err := queries.New(s.db).InsertUserToken(ctx, &queries.InsertUserTokenParams{
		UserID:     user.ID,
		Context:    "email_verification",
		Token:      crypto.HashUserToken(token),
		ValidUntil: time.Now().Add(validity),
	})

	if err != nil {
		return nil, fmt.Errorf("IssueEmailVerificationTokenForUser: %w", err)
	}

	return &types.UserTokenDTO{
		UserToken:      userToken,
		User:           user,
		PlaintextToken: token,
	}, nil
}
