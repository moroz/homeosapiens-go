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

type CartViewDto struct {
	CartItems  []*queries.GetCartItemsByCartIdRow
	GrandTotal decimal.Decimal
}

func NewCartService(db queries.DBTX) *CartService {
	return &CartService{db}
}

func (s *CartService) AddProductToCart(ctx context.Context, cartID *uuid.UUID, productID uuid.UUID) (*queries.CartLineItem, error) {
	if cartID == nil {
		cartID = new(uuid.Must(uuid.NewV7()))
	}

	item, err := queries.New(s.db).InsertCartLineItem(ctx, &queries.InsertCartLineItemParams{
		CartID:    *cartID,
		ProductID: productID,
	})
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (s *CartService) DeleteCartItem(ctx context.Context, cartID uuid.UUID, productID uuid.UUID) (bool, error) {
	id, err := queries.New(s.db).DeleteCartItem(ctx, &queries.DeleteCartItemParams{
		CartID:    cartID,
		ProductID: productID,
	})
	return id != uuid.UUID{}, err
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

func (c *CartViewDto) IsEmpty() bool {
	return len(c.CartItems) == 0
}
