package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sqlcrypter "github.com/bincyber/go-sqlcrypter"
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

// maybePlaintext decrypts an optional encrypted column, keeping SQL NULL as a
// JSON null rather than an empty string.
func maybePlaintext(b *sqlcrypter.EncryptedBytes) *string {
	if b == nil {
		return nil
	}
	return new(b.Plaintext())
}

// orderDetails maps an order and its line items onto the OrderDetails schema.
func orderDetails(o *queries.Order, items []*queries.OrderLineItem) OrderDetails {
	lineItems := make([]OrderLineItem, len(items))
	for i, it := range items {
		lineItems[i] = OrderLineItem{
			Id:                   it.ID,
			ProductId:            it.ProductID,
			ProductTitle:         it.ProductTitle,
			ProductPrice:         it.ProductPriceAmount.StringFixedBank(2),
			ProductPriceCurrency: it.ProductPriceCurrency,
			Quantity:             it.Quantity,
		}
	}

	return OrderDetails{
		Id:                      o.ID,
		UserId:                  o.UserID,
		OrderNumber:             fmt.Sprintf("%d", o.OrderNumber),
		Status:                  orderStatus(o),
		GrandTotal:              o.GrandTotal.StringFixedBank(2),
		Currency:                o.Currency,
		Email:                   openapi_types.Email(o.Email.Plaintext()),
		GivenName:               o.BillingGivenName.Plaintext(),
		FamilyName:              o.BillingFamilyName.Plaintext(),
		Phone:                   maybePlaintext(o.BillingPhone),
		AddressLine1:            o.BillingAddressLine1.Plaintext(),
		AddressLine2:            maybePlaintext(o.BillingAddressLine2),
		City:                    o.BillingCity.Plaintext(),
		PostalCode:              maybePlaintext(o.BillingPostalCode),
		BillingCountry:          o.BillingCountry,
		TaxId:                   maybePlaintext(o.BillingTaxID),
		PreferredLocale:         string(o.PreferredLocale),
		StripeCheckoutSessionId: o.StripeCheckoutSessionID,
		LineItems:               lineItems,
		InsertedAt:              o.InsertedAt,
		PaidAt:                  o.PaidAt,
		CancelledAt:             o.CancelledAt,
	}
}

func (s *orderServer) GetOrder(ctx context.Context, request GetOrderRequestObject) (GetOrderResponseObject, error) {
	o, err := s.q.GetOrderByID(ctx, request.Id)
	if errors.Is(err, sql.ErrNoRows) {
		return GetOrder404Response{}, nil
	}
	if err != nil {
		return nil, err
	}

	items, err := s.q.GetOrderLineItemsForOrderID(ctx, o.ID)
	if err != nil {
		return nil, err
	}

	return GetOrder200JSONResponse(orderDetails(o, items)), nil
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
