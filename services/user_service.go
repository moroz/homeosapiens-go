package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/alexedwards/argon2id"
	"github.com/moroz/homeosapiens-go/db/queries"
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
