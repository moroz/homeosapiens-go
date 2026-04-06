package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/bincyber/go-sqlcrypter"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/internal/crypto"
	"github.com/moroz/homeosapiens-go/internal/mailers"
	"github.com/moroz/homeosapiens-go/types"
)

type OrderService struct {
	db                   queries.DBTX
	paymentIntentService StripeService
	mailer               mailers.OrderMailer
}

func NewOrderService(db queries.DBTX, service StripeService, mailer mailers.OrderMailer) *OrderService {
	return &OrderService{
		db:                   db,
		paymentIntentService: service,
		mailer:               mailer,
	}
}

func maybeEncrypt(str string) *sqlcrypter.EncryptedBytes {
	if str == "" {
		return nil
	}
	return new(sqlcrypter.NewEncryptedBytes(str))
}

func (s *OrderService) CreateOrder(ctx context.Context, cartId uuid.UUID, user *queries.User, params *types.OrderParams) (*types.OrderDTO, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	db, ok := s.db.(*pgxpool.Pool)
	if !ok {
		return nil, fmt.Errorf("failed to cast connection as *pgxpool.Pool, got: %T", db)
	}

	tx, err := db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("CreateOrder: failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	cart, err := queries.New(tx).GetCart(ctx, cartId)
	if err != nil {
		return nil, fmt.Errorf("CreateOrder: %w", err)
	}

	if cart.ItemCount == 0 {
		return nil, fmt.Errorf("CreateOrder: cart is empty")
	}

	items, err := queries.New(tx).GetCartItemsByCartId(ctx, cartId)
	if err != nil {
		return nil, fmt.Errorf("CreateOrder: %w", err)
	}

	if user == nil {
		user, err = s.findOrCreateUserForOrder(ctx, tx, params)
		if err != nil {
			return nil, fmt.Errorf("CreateOrder: %w", err)
		}
	}

	var result types.OrderDTO

	result.Order, err = queries.New(tx).InsertOrder(ctx, &queries.InsertOrderParams{
		PreferredLocale:     queries.Locale(params.PreferredLocale),
		UserID:              user.ID,
		GrandTotal:          cart.ProductTotal,
		Currency:            "PLN",
		BillingGivenName:    sqlcrypter.NewEncryptedBytes(params.BillingGivenName),
		BillingFamilyName:   sqlcrypter.NewEncryptedBytes(params.BillingFamilyName),
		BillingPhone:        maybeEncrypt(params.BillingPhone),
		BillingAddressLine1: sqlcrypter.NewEncryptedBytes(params.BillingAddressLine1),
		BillingAddressLine2: maybeEncrypt(params.BillingAddressLine1),
		BillingCity:         sqlcrypter.NewEncryptedBytes(params.BillingCity),
		BillingPostalCode:   maybeEncrypt(params.BillingPostalCode),
		BillingCountry:      params.BillingCountry,
		Email:               sqlcrypter.NewEncryptedBytes(params.Email),
		BillingTaxID:        maybeEncrypt(params.BillingTaxID),
	})
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		if item.Event.BasePriceAmount == nil {
			return nil, fmt.Errorf("CreateOrder: event %s has no price set", item.Event.ID)
		}

		lineItem, err := queries.New(tx).InsertOrderLineItem(ctx, &queries.InsertOrderLineItemParams{
			OrderID:          result.Order.ID,
			EventID:          item.Event.ID,
			EventTitle:       item.Event.TitleEn,
			EventPriceAmount: item.Subtotal,
		})

		if err != nil {
			return nil, fmt.Errorf("CreateOrder: %w", err)
		}

		result.LineItems = append(result.LineItems, lineItem)
	}

	if err := queries.New(tx).DeleteCart(ctx, cartId); err != nil {
		return nil, fmt.Errorf("CreateOrder: %w", err)
	}

	result.CheckoutSession, err = s.paymentIntentService.CreateCheckoutSession(ctx, &result)
	if err != nil {
		return nil, fmt.Errorf("CreateOrder: stripe error: %w", err)
	}

	result.Order, err = queries.New(tx).StoreCheckoutSessionIDOnOrder(ctx, &queries.StoreCheckoutSessionIDOnOrderParams{
		StripeCheckoutSessionID: &result.CheckoutSession.ID,
		ID:                      result.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("CreateOrder: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("CreateOrder: could not commit transaction: %w", err)
	}

	if err := s.mailer.SendOrderConfirmation(ctx, &result); err != nil {
		log.Printf("Error sending order confirmation email for order %s: %s", result.Order.ID, err)
	}

	return &result, nil
}

func (s *OrderService) findOrCreateUserForOrder(ctx context.Context, tx pgx.Tx, params *types.OrderParams) (*queries.User, error) {
	user, err := NewUserService(tx).FindUserByEmail(ctx, params.Email)
	if err == nil {
		return user, nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		// TODO: Should this be in UserService?
		return queries.New(tx).InsertUser(ctx, &queries.InsertUserParams{
			PreferredLocale: queries.Locale(params.PreferredLocale),
			Email:           sqlcrypter.NewEncryptedBytes(params.Email),
			EmailHash:       crypto.HashEmail(params.Email),
			GivenName:       sqlcrypter.NewEncryptedBytes(params.BillingGivenName),
			FamilyName:      sqlcrypter.NewEncryptedBytes(params.BillingFamilyName),
			Country:         &params.BillingCountry,
		})
	}

	return nil, err
}

func (s *OrderService) GetOrderByCheckoutSessionID(ctx context.Context, sessionID string) (*queries.Order, error) {
	return queries.New(s.db).GetOrderByCheckoutSessionID(ctx, sessionID)
}

var ErrOrderAlreadyPaid = fmt.Errorf("MarkOrderPaidByCheckoutSessionID: order has already been marked as paid")

func (s *OrderService) MarkOrderPaidByCheckoutSessionID(ctx context.Context, sessionID string) (*queries.Order, error) {
	conn, ok := s.db.(*pgxpool.Pool)
	if !ok {
		return nil, fmt.Errorf("want *pgxpool.Conn, got %T", s.db)
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// Locks the row with SELECT FOR UPDATE
	order, err := queries.New(tx).GetOrderByCheckoutSessionIDForUpdate(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("MarkOrderPaidByCheckoutSessionID: %w", err)
	}

	if order.PaidAt != nil {
		return nil, ErrOrderAlreadyPaid
	}

	order, err = queries.New(tx).MarkOrderAsPaid(ctx, order.ID)
	if err != nil {
		return nil, fmt.Errorf("MarkOrderPaidByCheckoutSessionID: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("MarkOrderPaidByCheckoutSessionID: %w", err)
	}

	return order, err
}

func (s *OrderService) GetOrderDetails(ctx context.Context, orderID uuid.UUID) (*types.OrderDTO, error) {
	var result types.OrderDTO
	var err error

	result.Order, err = queries.New(s.db).GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	result.LineItems, err = queries.New(s.db).GetOrderLineItemsForOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
