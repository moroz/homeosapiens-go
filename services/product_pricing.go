package services

import (
	"context"

	"uuid"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jackc/pgx/v5"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/shopspring/decimal"
)

// patchProductPricingParams describes a pricing change applied to whatever a
// product hangs off — an event or a video group.
type patchProductPricingParams struct {
	// ProductID is the product the record currently points at, nil when free.
	ProductID *uuid.UUID
	Price     types.Optional[decimal.Decimal]
	Currency  types.Optional[string]

	// ProductType, TitlePl and TitleEn describe the product to create if the
	// record is being given a price for the first time.
	ProductType queries.ProductType
	TitlePl     string
	TitleEn     string
}

// patchProductPricing applies a pricing change and returns the product ID the
// record must be pointed at, or nil when the existing association already holds.
// Events and video groups share it so that the "a record that already has a
// product keeps it" rule cannot drift apart between the two.
func patchProductPricing(ctx context.Context, tx pgx.Tx, params *patchProductPricingParams) (*uuid.UUID, error) {
	if !params.Price.Set && !params.Currency.Set {
		return nil, nil
	}

	// A record that already has a product keeps it, even when the price drops to
	// zero: cart and order line items reference it by ID.
	if params.ProductID != nil {
		product, err := queries.New(tx).GetProductById(ctx, *params.ProductID)
		if err != nil {
			return nil, err
		}

		amount := product.BasePriceAmount
		if params.Price.Set {
			amount = params.Price.Value
		}

		currency := product.BasePriceCurrency
		if params.Currency.Set {
			currency = params.Currency.Value
		}

		_, err = queries.New(tx).UpdateProductPrice(ctx, &queries.UpdateProductPriceParams{
			ProductID:         product.ID,
			BasePriceAmount:   amount,
			BasePriceCurrency: currency,
		})

		return nil, err
	}

	// A free record stays free until it is given a non-zero price; a currency on
	// its own has nowhere to be stored.
	if !params.Price.Set || params.Price.Value.Equal(decimal.Zero) {
		return nil, nil
	}

	if !params.Currency.Set {
		return nil, validation.Errors{
			"currency": validation.NewError("required", "is required when setting a price"),
		}
	}

	product, err := queries.New(tx).InsertProduct(ctx, &queries.InsertProductParams{
		ProductType:       params.ProductType,
		TitlePl:           params.TitlePl,
		TitleEn:           params.TitleEn,
		BasePriceAmount:   params.Price.Value,
		BasePriceCurrency: params.Currency.Value,
	})
	if err != nil {
		return nil, err
	}

	return &product.ID, nil
}
