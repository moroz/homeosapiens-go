package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/types"
)

type blogPostServer struct {
	db *pgxpool.Pool
}

func NewBlogPostServer(db *pgxpool.Pool) *blogPostServer {
	return &blogPostServer{db: db}
}

func (s *blogPostServer) ListBlogPosts(ctx context.Context, _ ListBlogPostsRequestObject) (ListBlogPostsResponseObject, error) {
	posts, err := queries.New(s.db).ListAllBlogPosts(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]BlogPost, len(posts))
	for i, row := range posts {
		result[i] = BlogPost{
			Id:          row.ID,
			Body:        row.Body,
			Slug:        row.Slug,
			PublishedAt: row.PublishedAt,
			InsertedAt:  row.InsertedAt,
			Language:    string(row.Language),
			Title:       row.Title,
			UpdatedAt:   row.UpdatedAt,
		}
	}
	return ListBlogPosts200JSONResponse(result), nil
}

func (s *blogPostServer) CreateBlogPost(ctx context.Context, request CreateBlogPostRequestObject) (CreateBlogPostResponseObject, error) {
	p := request.Body

	post, err := services.NewBlogService(s.db).CreateBlogPost(ctx, &types.CreateBlogPostInput{
		Title:    p.Title,
		Slug:     p.Slug,
		Language: p.Language,
		Body:     p.Body,
	})
	if verr, ok := errors.AsType[validation.Errors](err); ok {
		return CreateBlogPost422JSONResponse{Errors: validationErrorMessages(verr)}, nil
	}
	if err != nil {
		return nil, err
	}

	return CreateBlogPost201JSONResponse{
		Body: BlogPost{
			Body:        post.Body,
			Id:          post.ID,
			InsertedAt:  post.InsertedAt,
			Language:    string(post.Language),
			PublishedAt: post.PublishedAt,
			Slug:        post.Slug,
			Title:       post.Title,
			UpdatedAt:   post.UpdatedAt,
		},
		Headers: CreateBlogPost201ResponseHeaders{
			Location: new(fmt.Sprintf("/api/admin/blog-posts/%s", post.ID)),
		},
	}, nil
}

func (s *blogPostServer) GetBlogPost(ctx context.Context, request GetBlogPostRequestObject) (GetBlogPostResponseObject, error) {
	post, err := queries.New(s.db).GetBlogPostById(ctx, request.Id)
	if errors.Is(err, sql.ErrNoRows) {
		return GetBlogPost404Response{}, nil
	}
	if err != nil {
		return nil, err
	}

	return GetBlogPost200JSONResponse{
		Body:        post.Body,
		Id:          post.ID,
		InsertedAt:  post.InsertedAt,
		Language:    string(post.Language),
		PublishedAt: post.PublishedAt,
		Slug:        post.Slug,
		Title:       post.Title,
		UpdatedAt:   post.UpdatedAt,
	}, nil
}

func (s *blogPostServer) UpdateBlogPost(ctx context.Context, request UpdateBlogPostRequestObject) (UpdateBlogPostResponseObject, error) {
	p := request.Body

	post, err := services.NewBlogService(s.db).UpdateBlogPost(ctx, request.Id, &types.UpdateBlogPostInput{
		Title:    p.Title,
		Slug:     p.Slug,
		Language: p.Language,
		Body:     p.Body,
	})
	if errors.Is(err, sql.ErrNoRows) {
		return UpdateBlogPost404Response{}, nil
	}
	if verr, ok := errors.AsType[validation.Errors](err); ok {
		return UpdateBlogPost422JSONResponse{Errors: validationErrorMessages(verr)}, nil
	}
	if err != nil {
		return nil, err
	}

	return UpdateBlogPost200JSONResponse{
		Body:        post.Body,
		Id:          post.ID,
		InsertedAt:  post.InsertedAt,
		Language:    string(post.Language),
		PublishedAt: post.PublishedAt,
		Slug:        post.Slug,
		Title:       post.Title,
		UpdatedAt:   post.UpdatedAt,
	}, nil
}

func (s *blogPostServer) PublishBlogPost(ctx context.Context, request PublishBlogPostRequestObject) (PublishBlogPostResponseObject, error) {
	_, err := services.NewBlogService(s.db).PublishBlogPost(ctx, request.Id)
	if errors.Is(err, sql.ErrNoRows) {
		return PublishBlogPost404Response{}, nil
	}
	if verr, ok := errors.AsType[validation.Errors](err); ok {
		return PublishBlogPost422JSONResponse{Errors: validationErrorMessages(verr)}, nil
	}
	if err != nil {
		return nil, err
	}

	return PublishBlogPost204Response{}, nil
}

func (s *blogPostServer) UnpublishBlogPost(ctx context.Context, request UnpublishBlogPostRequestObject) (UnpublishBlogPostResponseObject, error) {
	_, err := services.NewBlogService(s.db).UnpublishBlogPost(ctx, request.Id)
	if errors.Is(err, sql.ErrNoRows) {
		return UnpublishBlogPost404Response{}, nil
	}
	if verr, ok := errors.AsType[validation.Errors](err); ok {
		return UnpublishBlogPost422JSONResponse{Errors: validationErrorMessages(verr)}, nil
	}
	if err != nil {
		return nil, err
	}

	return UnpublishBlogPost204Response{}, nil
}
