package api

import (
	"context"

	"github.com/moroz/homeosapiens-go/db/queries"
)

type productServer struct {
	q *queries.Queries
}

func NewProductServer(db queries.DBTX) *productServer {
	return &productServer{q: queries.New(db)}
}

// product maps a database row onto the Product schema.
func product(p *queries.Product) Product {
	return Product{
		Id:                p.ID,
		ProductType:       p.ProductType,
		TitlePl:           p.TitlePl,
		TitleEn:           p.TitleEn,
		BasePrice:         p.BasePriceAmount.StringFixedBank(2),
		BasePriceCurrency: p.BasePriceCurrency,
		InsertedAt:        p.InsertedAt,
		UpdatedAt:         p.UpdatedAt,
	}
}

func (s *productServer) ListProducts(ctx context.Context, request ListProductsRequestObject) (ListProductsResponseObject, error) {
	page, perPage := resolvePaginationParams(request.Params.Page, request.Params.PerPage)

	products, err := s.q.PaginateProducts(ctx, &queries.PaginateProductsParams{
		Page:    page,
		PerPage: perPage,
	})
	if err != nil {
		return nil, err
	}

	count, err := s.q.CountProducts(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]Product, len(products))
	for i, p := range products {
		out[i] = product(p)
	}

	return ListProducts200JSONResponse{
		Data: out,
		Pagination: Pagination{
			Page:       page,
			PerPage:    perPage,
			Total:      count,
			TotalPages: countPages(count, perPage),
		},
	}, nil
}
