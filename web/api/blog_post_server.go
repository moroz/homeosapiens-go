package api

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moroz/homeosapiens-go/db/queries"
)

type blogPostServer struct {
	db *pgxpool.Pool
}

func NewBlogPostServer(db *pgxpool.Pool) *blogPostServer {
	return &blogPostServer{db: db}
}

func (s *blogPostServer) ListBlogPosts(ctx context.Context, request ListBlogPostsRequestObject) (ListBlogPostsResponseObject, error) {
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
