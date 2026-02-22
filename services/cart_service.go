package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/shopspring/decimal"
)

type CartService struct {
	db queries.DBTX
}

func NewCartService(db queries.DBTX) *CartService {
	return &CartService{db}
}

func (s *CartService) AddEventToCart(ctx context.Context, cartId *uuid.UUID, eventId uuid.UUID) (*queries.CartLineItem, error) {
	if cartId == nil {
		cartId = new(uuid.Must(uuid.NewV7()))
	}

	item, err := queries.New(s.db).InsertCartLineItem(ctx, &queries.InsertCartLineItemParams{
		CartID:  *cartId,
		EventID: eventId,
	})
	if err != nil {
		return nil, err
	}

	return item, nil
}

type CartViewDto struct {
	CartItems  []*queries.GetCartItemsByCartIdRow
	GrandTotal decimal.Decimal
}

func (s *CartService) GetCartItemsByCartId(ctx context.Context, cartId *uuid.UUID) (*CartViewDto, error) {
	result := CartViewDto{
		CartItems:  nil,
		GrandTotal: decimal.Zero,
	}

	if cartId == nil {
		return &result, nil
	}

	items, err := queries.New(s.db).GetCartItemsByCartId(ctx, *cartId)
	if err != nil {
		return &result, err
	}
	result.CartItems = items

	for _, item := range items {
		result.GrandTotal = result.GrandTotal.Add(item.Subtotal)
	}

	return &result, nil
}
