package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/types"
)

type OrderRepository struct {
	db queries.DBTX
}

func NewOrderRepository(db queries.DBTX) *OrderRepository {
	return &OrderRepository{db}
}

func (s *OrderRepository) GetOrderDetails(ctx context.Context, orderID uuid.UUID) (*types.OrderDTO, error) {
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

func (s *OrderRepository) MarkOrderPaidByCheckoutSessionID(ctx context.Context, sessionID string) (*queries.Order, error) {
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
