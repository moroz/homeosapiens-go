package services

import (
	"context"

	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/types"
)

type BlogService struct {
	db queries.DBTX
}

func NewBlogService(db queries.DBTX) *BlogService {
	return &BlogService{db}
}

func (s *BlogService) CreateBlogPost(ctx context.Context, params *types.CreateBlogPostInput) (*queries.BlogPost, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	return queries.New(s.db).InsertBlogPost(ctx, &queries.InsertBlogPostParams{
		Title:    params.Title,
		Slug:     params.Slug,
		Body:     params.Body,
		Language: queries.Locale(params.Language),
	})
}
