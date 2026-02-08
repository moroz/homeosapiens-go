package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/alexedwards/argon2id"
	"github.com/bincyber/go-sqlcrypter"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/internal/crypto"
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
	normalizedEmail := crypto.HashEmail(email)
	user, err := queries.New(s.db).GetUserByEmail(ctx, normalizedEmail)
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

func (s *UserService) CreateUser(ctx context.Context, params *types.CreateUserParams) (*queries.User, error) {
	passwordHash, err := argon2id.CreateHash(params.Password, config.ResolveArgon2Params())
	if err != nil {
		return nil, err
	}

	return queries.New(s.db).UpsertUserFromSeedData(ctx, &queries.UpsertUserFromSeedDataParams{
		Email:        sqlcrypter.NewEncryptedBytes(params.Email),
		EmailHash:    crypto.HashEmail(params.Email),
		GivenName:    sqlcrypter.NewEncryptedBytes(params.GivenName),
		FamilyName:   sqlcrypter.NewEncryptedBytes(params.FamilyName),
		Country:      &params.Country,
		PasswordHash: &passwordHash,
		UserRole:     params.Role,
	})
}

func (s *UserService) FindOrCreateUserFromClaims(ctx context.Context, claims *types.GoogleIDTokenClaims) (*queries.User, error) {
	return queries.New(s.db).FindOrCreateUserFromClaims(ctx, &queries.FindOrCreateUserFromClaimsParams{
		Email:          sqlcrypter.NewEncryptedBytes(claims.Email),
		EmailHash:      crypto.HashEmail(claims.Email),
		GivenName:      sqlcrypter.NewEncryptedBytes(claims.GivenName),
		FamilyName:     sqlcrypter.NewEncryptedBytes(claims.FamilyName),
		ProfilePicture: &claims.Avatar,
		EmailConfirmed: true,
	})
}

func (s *UserService) FindOrCreateUserFromEventRegistrationParams(ctx context.Context, params *types.CreateEventRegistrationParams) (*queries.User, error) {
	return queries.New(s.db).FindOrCreateUserFromClaims(ctx, &queries.FindOrCreateUserFromClaimsParams{
		Email:          sqlcrypter.NewEncryptedBytes(params.Email),
		EmailHash:      crypto.HashEmail(params.Email),
		GivenName:      sqlcrypter.NewEncryptedBytes(params.GivenName),
		FamilyName:     sqlcrypter.NewEncryptedBytes(params.FamilyName),
		EmailConfirmed: false,
	})
}

func (s *UserService) UpdateUserProfile(ctx context.Context, user *queries.User, params *types.UpdateProfileRequest) (*queries.User, error) {
	var profession *string
	if strings.TrimSpace(params.Profession) != "" {
		profession = &params.Profession
	}

	var licenceNumber *sqlcrypter.EncryptedBytes
	if strings.TrimSpace(params.LicenceNumber) != "" {
		encrypted := sqlcrypter.NewEncryptedBytes(params.LicenceNumber)
		licenceNumber = &encrypted
	}

	return queries.New(s.db).UpdateUserProfile(ctx, &queries.UpdateUserProfileParams{
		GivenName:     sqlcrypter.NewEncryptedBytes(params.GivenName),
		FamilyName:    sqlcrypter.NewEncryptedBytes(params.FamilyName),
		Profession:    profession,
		LicenceNumber: licenceNumber,
		Country:       &params.Country,
		ID:            user.ID,
	})
}
