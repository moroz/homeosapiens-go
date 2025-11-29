package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/alexedwards/argon2id"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/types"
)

type UserService struct {
	db queries.DBTX
}

func NewUserService(db queries.DBTX) *UserService {
	return &UserService{db}
}

var ErrInvalidPassword = errors.New("invalid password")
var ErrNoPasswordHash = errors.New("user has no password set")

func (s *UserService) AuthenticateUserByEmailPassword(ctx context.Context, email, password string) (*queries.User, error) {
	tmpl := "AuthenticateUserByEmailPassword: %w"
	user, err := queries.New(s.db).GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf(tmpl, err)
	}

	if user.PasswordHash == nil {
		return nil, fmt.Errorf(tmpl, ErrNoPasswordHash)
	}

	var ok bool
	if ok, _, err = argon2id.CheckHash(password, *user.PasswordHash); !ok {
		return nil, fmt.Errorf(tmpl, ErrInvalidPassword)
	} else if err != nil {
		return nil, fmt.Errorf(tmpl, err)
	}
	return user, nil
}

func (s *UserService) FindOrCreateUserFromClaims(ctx context.Context, claims *types.GoogleIDTokenClaims) (*queries.User, error) {
	return queries.New(s.db).FindOrCreateUserFromClaims(ctx, &queries.FindOrCreateUserFromClaimsParams{
		Email:          claims.Email,
		GivenName:      claims.GivenName,
		FamilyName:     claims.FamilyName,
		ProfilePicture: &claims.Avatar,
	})
}
