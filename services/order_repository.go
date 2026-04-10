package services

import (
	"context"

	"github.com/google/uuid"
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

func (s *OrderRepository) GetOrderDetailsByCheckoutSessionID(ctx context.Context, sessionID string) (*types.OrderDTO, error) {
	var result types.OrderDTO
	var err error

	result.Order, err = queries.New(s.db).GetOrderByCheckoutSessionID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	result.LineItems, err = queries.New(s.db).GetOrderLineItemsForOrderID(ctx, result.Order.ID)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
