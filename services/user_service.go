package services

import (
	"context"
	"errors"
	"fmt"
	"net/netip"
	"strings"

	"github.com/alexedwards/argon2id"
	"github.com/bincyber/go-sqlcrypter"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
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
	user, err := s.FindUserByEmail(ctx, email)
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

func (s *UserService) FindUserByEmail(ctx context.Context, email string) (*queries.User, error) {
	normalizedEmail := crypto.HashEmail(email)
	return queries.New(s.db).GetUserByEmail(ctx, normalizedEmail)
}

func (s *UserService) CreateUser(ctx context.Context, params *types.SeedUserParams) (*queries.User, error) {
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

func (s *UserService) FindOrCreateUserFromClaims(ctx context.Context, locale string, claims *types.GoogleIDTokenClaims) (*queries.User, error) {

	return queries.New(s.db).FindOrCreateUserFromClaims(ctx, &queries.FindOrCreateUserFromClaimsParams{
		PreferredLocale: queries.Locale(locale),
		Email:           sqlcrypter.NewEncryptedBytes(claims.Email),
		EmailHash:       crypto.HashEmail(claims.Email),
		GivenName:       sqlcrypter.NewEncryptedBytes(claims.GivenName),
		FamilyName:      sqlcrypter.NewEncryptedBytes(claims.FamilyName),
		ProfilePicture:  &claims.Avatar,
		EmailConfirmed:  true,
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

func (s *UserService) ListUsers(ctx context.Context) ([]*queries.User, error) {
	return queries.New(s.db).ListUsers(ctx)
}

func (s *UserService) SetUserLastLogin(ctx context.Context, ipAddr string, userId uuid.UUID) error {
	parsed, err := netip.ParseAddr(ipAddr)
	if err != nil {
		return err
	}

	return queries.New(s.db).SetUserLastLogin(ctx, &queries.SetUserLastLoginParams{
		LastLoginIp: &parsed,
		ID:          userId,
	})
}

func (s *UserService) RegisterUser(ctx context.Context, params *types.RegisterUserParams) (*queries.User, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	passwordHash, err := argon2id.CreateHash(params.Password, config.ResolveArgon2Params())
	if err != nil {
		return nil, err
	}

	tx, err := s.db.(*pgxpool.Pool).Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("RegisterUser: %w", err)
	}
	defer tx.Rollback(ctx)

	user, err := queries.New(tx).InsertUser(ctx, &queries.InsertUserParams{
		Email:        sqlcrypter.NewEncryptedBytes(params.Email),
		EmailHash:    crypto.HashEmail(params.Email),
		GivenName:    sqlcrypter.NewEncryptedBytes(params.GivenName),
		FamilyName:   sqlcrypter.NewEncryptedBytes(params.FamilyName),
		PasswordHash: &passwordHash,
	})

	if err, ok := errors.AsType[*pgconn.PgError](err); ok && err.Code == "23505" && err.ConstraintName == "users_email_hash_key" {
		return nil, validation.Errors{
			"email": validation.NewError("unique", "has already been taken"),
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("RegisterUser: %w", err)
	}

	return user, nil
}
