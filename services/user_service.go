package services

import (
	"context"
	"errors"

	"github.com/alexedwards/argon2id"
	"github.com/moroz/homeosapiens-go/db/queries"
)

type UserService struct {
	db queries.DBTX
}

func NewUserService(db queries.DBTX) *UserService {
	return &UserService{db}
}

func (s *UserService) AuthenticateUserByEmailPassword(ctx context.Context, email, password string) (*queries.User, error) {
	user, err := queries.New(s.db).GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if user.PasswordHash == nil {
		return nil, errors.New("user has no password set")
	}

	ok, _, err := argon2id.CheckHash(password, *user.PasswordHash)
	if ok {
		return user, nil
	}

	return nil, err
}
