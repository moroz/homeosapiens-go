package services

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moroz/homeosapiens-go/db/queries"
)

type CartService struct {
	db queries.DBTX
}

func NewCartItemService(db queries.DBTX) *CartService {
	return &CartService{db}
}

func (s *CartService) AddEventToCart(ctx context.Context, cart *queries.Cart, user *queries.User, eventId pgtype.UUID) (*queries.CartLineItem, error) {
	conn, ok := s.db.(*pgxpool.Conn)
	if !ok {
		return nil, fmt.Errorf("failed to cast db as *pgxpool.Conn: %T", s.db)
	}
	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	if cart == nil {
		userId := pgtype.UUID{}
		if user != nil {
			userId = user.ID
		}

		cart, err = queries.New(tx).InsertCart(ctx, userId)
		if err != nil {
			return nil, err
		}
	}

	item, err := queries.New(tx).InsertCartLineItem(ctx, &queries.InsertCartLineItemParams{
		CartID:  cart.ID,
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
