package services

import (
	"context"

	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/types"
)

type UserPasswordResetService struct {
	db queries.DBTX
}

func NewUserPasswordResetService(db queries.DBTX) *UserPasswordResetService {
	return &UserPasswordResetService{db}
}

func (s *UserPasswordResetService) IssuePasswordResetTokenForUser(ctx context.Context, user *queries.User) (*types.UserTokenDTO, error) {
	return NewUserTokenService(s.db).IssueHashedTokenForUser(ctx, user, "password_reset", config.PasswordResetTokenValidity)
}
