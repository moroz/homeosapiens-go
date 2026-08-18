package services

import (
	"context"
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
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

	post, err := queries.New(s.db).InsertBlogPost(ctx, &queries.InsertBlogPostParams{
		Title:    params.Title,
		Slug:     params.Slug,
		Body:     params.Body,
		Language: queries.Locale(params.Language),
	})
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok && pgErr.Code == "23505" && pgErr.ConstraintName == "blog_posts_slug_key" {
		return nil, validation.Errors{
			"slug": validation.NewError("unique", "has already been taken"),
		}
	}
	return post, err
}

func (s *BlogService) UpdateBlogPost(ctx context.Context, id uuid.UUID, params *types.UpdateBlogPostInput) (*queries.BlogPost, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	post, err := queries.New(s.db).UpdateBlogPost(ctx, &queries.UpdateBlogPostParams{
		ID:       id,
		Title:    params.Title,
		Slug:     params.Slug,
		Body:     params.Body,
		Language: queries.Locale(params.Language),
	})
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok && pgErr.Code == "23505" && pgErr.ConstraintName == "blog_posts_slug_key" {
		return nil, validation.Errors{
			"slug": validation.NewError("unique", "has already been taken"),
		}
	}
	return post, err
}

func (s *BlogService) PublishBlogPost(ctx context.Context, id uuid.UUID) (*queries.BlogPost, error) {
	post, err := queries.New(s.db).GetBlogPostById(ctx, id)
	if err != nil {
		return nil, err
	}

	params := &types.PublishBlogPostValidation{
		Body:        post.Body,
		PublishedAt: post.PublishedAt,
	}
	if err := params.Validate(); err != nil {
		return nil, err
	}

	return queries.New(s.db).PublishBlogPost(ctx, id)
}

func (s *BlogService) UnpublishBlogPost(ctx context.Context, id uuid.UUID) (*queries.BlogPost, error) {
	post, err := queries.New(s.db).GetBlogPostById(ctx, id)
	if err != nil {
		return nil, err
	}

	params := &types.UnpublishBlogPostValidation{
		PublishedAt: post.PublishedAt,
	}
	if err := params.Validate(); err != nil {
		return nil, err
	}

	return queries.New(s.db).UnpublishBlogPost(ctx, id)
}
