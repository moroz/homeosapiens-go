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

func (s *UserTokenService) IssueHashedTokenForUser(ctx context.Context, user *queries.User, tokenContext string, validity time.Duration) (*types.UserTokenDTO, error) {
	token, err := generateToken()
	if err != nil {
		return nil, fmt.Errorf("IssueHashedTokenForUser: %w", err)
	}

	userToken, err := queries.New(s.db).InsertUserToken(ctx, &queries.InsertUserTokenParams{
		UserID:     user.ID,
		Context:    tokenContext,
		Token:      crypto.HashUserToken(token),
		ValidUntil: time.Now().Add(validity),
	})

	if err != nil {
		return nil, fmt.Errorf("IssueHashedTokenForUser: %w", err)
	}

	return &types.UserTokenDTO{
		UserToken:      userToken,
		User:           user,
		PlaintextToken: token,
	}, nil
}

// IssueEmailVerificationTokenForUser issues a hashed `UserToken` with context set to `email_verification`. These tokens are used in the account activation flow, and are encoded in their plaintext form in email verification URLs.
func (s *UserTokenService) IssueEmailVerificationTokenForUser(ctx context.Context, user *queries.User, validity time.Duration) (*types.UserTokenDTO, error) {
	return s.IssueHashedTokenForUser(ctx, user, "email_verification", validity)
}

// IssueUserRegistrationTokenForUser issues a hashed `UserToken` with context set to `user_registration`. These tokens are used to display the user registration success page with relevant information (such as the user's email address and a link to resend the account validation email).
func (s *UserTokenService) IssueUserRegistrationTokenForUser(ctx context.Context, user *queries.User, validity time.Duration) (*types.UserTokenDTO, error) {
	return s.IssueHashedTokenForUser(ctx, user, "user_registration", validity)
}
