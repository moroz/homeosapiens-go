package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moroz/homeosapiens-go/db/queries"
)

type CartService struct {
	db queries.DBTX
}

func NewCartItemService(db queries.DBTX) *CartService {
	return &CartService{db}
}

func (s *CartService) AddEventToCart(ctx context.Context, cartId *uuid.UUID, user *queries.User, eventId uuid.UUID) (*queries.CartLineItem, error) {
	conn, ok := s.db.(*pgxpool.Pool)
	if !ok {
		return nil, fmt.Errorf("failed to cast db as *pgxpool.Conn: %T", s.db)
	}
	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	if cartId == nil {
		var userId *uuid.UUID
		if user != nil {
			userId = &user.ID
		}

		cart, err := queries.New(tx).InsertCart(ctx, userId)
		if err != nil {
			return nil, err
		}

		cartId = &cart.ID
	}

	item, err := queries.New(tx).InsertCartLineItem(ctx, &queries.InsertCartLineItemParams{
		CartID:  *cartId,
		EventID: eventId,
	})
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return item, nil
}
