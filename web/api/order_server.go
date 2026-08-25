package api

import (
	"context"
	"fmt"

	"github.com/moroz/homeosapiens-go/db/queries"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type orderServer struct {
	q *queries.Queries
}

func NewOrderServer(db queries.DBTX) *orderServer {
	return &orderServer{q: queries.New(db)}
}

// orderStatus derives the status shown in the admin from the order timestamps.
// A cancelled order stays cancelled even if it had been paid before.
func orderStatus(o *queries.Order) OrderStatus {
	switch {
	case o.CancelledAt != nil:
		return Cancelled
	case o.PaidAt != nil:
		return Paid
	default:
		return Pending
	}
}

// order maps an order row onto the Order schema, decrypting the billing PII.
func order(o *queries.Order) Order {
	return Order{
		Id:             o.ID,
		OrderNumber:    fmt.Sprintf("%d", o.OrderNumber),
		Status:         orderStatus(o),
		GrandTotal:     o.GrandTotal.StringFixedBank(2),
		Currency:       o.Currency,
		Email:          openapi_types.Email(o.Email.Plaintext()),
		GivenName:      o.BillingGivenName.Plaintext(),
		FamilyName:     o.BillingFamilyName.Plaintext(),
		BillingCountry: o.BillingCountry,
		InsertedAt:     o.InsertedAt,
		PaidAt:         o.PaidAt,
		CancelledAt:    o.CancelledAt,
	}
}

func (s *orderServer) ListOrders(ctx context.Context, request ListOrdersRequestObject) (ListOrdersResponseObject, error) {
	page, perPage := resolvePaginationParams(request.Params.Page, request.Params.PerPage)

	orders, err := s.q.PaginateOrders(ctx, &queries.PaginateOrdersParams{
		Page:    page,
		PerPage: perPage,
	})
	if err != nil {
		return nil, err
	}

	count, err := s.q.CountOrders(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]Order, len(orders))
	for i, o := range orders {
		out[i] = order(o)
	}

	return ListOrders200JSONResponse{
		Data: out,
		Pagination: Pagination{
			Page:       page,
			PerPage:    perPage,
			Total:      count,
			TotalPages: countPages(count, perPage),
		},
	}, nil
}
