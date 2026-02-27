package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/bincyber/go-sqlcrypter"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/internal/crypto"
	"github.com/moroz/homeosapiens-go/types"
)

type OrderService struct {
	db queries.DBTX
}

func NewOrderService(db queries.DBTX) *OrderService {
	return &OrderService{
		db: db,
	}
}

func maybeEncrypt(str *string) *sqlcrypter.EncryptedBytes {
	if str == nil {
		return nil
	}
	return new(sqlcrypter.NewEncryptedBytes(*str))
}

func (s *OrderService) CreateOrder(ctx context.Context, cartId uuid.UUID, user *queries.User, params *types.OrderParams) (*queries.Order, error) {
	db, ok := s.db.(*pgxpool.Pool)
	if !ok {
		return nil, fmt.Errorf("failed to cast connection as *pgxpool.Pool, got: %T", db)
	}

	tx, err := db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	items, err := queries.New(tx).GetCart(ctx, cartId)
	if err != nil {
		return nil, err
	}

	if user == nil {
		user, err = s.findOrCreateUserForOrder(ctx, tx, params)
		if err != nil {
			return nil, err
		}
	}

	order, err := queries.New(tx).InsertOrder(ctx, &queries.InsertOrderParams{
		UserID:                 user.ID,
		GrandTotal:             items.ProductTotal,
		Currency:               "PLN",
		BillingGivenName:       sqlcrypter.NewEncryptedBytes(params.BillingGivenName),
		BillingFamilyName:      sqlcrypter.NewEncryptedBytes(params.BillingFamilyName),
		BillingPhone:           maybeEncrypt(params.BillingPhone),
		BillingStreet:          sqlcrypter.NewEncryptedBytes(params.BillingStreet),
		BillingHouseNumber:     sqlcrypter.NewEncryptedBytes(params.BillingHouseNumber),
		BillingApartmentNumber: maybeEncrypt(params.BillingApartmentNumber),
		BillingCity:            sqlcrypter.NewEncryptedBytes(params.BillingCity),
		BillingPostalCode:      maybeEncrypt(params.BillingPostalCode),
		BillingCountry:         params.BillingCountry,
		Email:                  sqlcrypter.NewEncryptedBytes(params.Email),
	})
	if err != nil {
		return nil, err
	}

	if err := queries.New(tx).DeleteCart(ctx, cartId); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) findOrCreateUserForOrder(ctx context.Context, tx pgx.Tx, params *types.OrderParams) (*queries.User, error) {
	user, err := NewUserService(tx).FindUserByEmail(ctx, params.Email)
	if user != nil {
		return user, nil
	}

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		// TODO: Should this be in UserService?
		return queries.New(tx).InsertUser(ctx, &queries.InsertUserParams{
			Email:      sqlcrypter.NewEncryptedBytes(params.Email),
			EmailHash:  crypto.HashEmail(params.Email),
			GivenName:  sqlcrypter.NewEncryptedBytes(params.BillingGivenName),
			FamilyName: sqlcrypter.NewEncryptedBytes(params.BillingFamilyName),
			Country:    &params.BillingCountry,
		})
	}

	return nil, err
}
