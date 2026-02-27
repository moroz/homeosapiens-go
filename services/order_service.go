package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/types"
)

type OrderService struct {
	db queries.DBTX
}

func (s *OrderService) CreateOrder(ctx context.Context, cartId uuid.UUID, user *queries.User, params *types.OrderParams) (*queries.Order, error) {

}
