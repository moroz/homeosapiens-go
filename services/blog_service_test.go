package services_test

import (
	"database/sql"
	"errors"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/services/mocks"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlogService_CreateBlogPost(t *testing.T) {
	ctx := t.Context()
	db, err := initDB(ctx)
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Exec(ctx, "truncate blog_posts cascade")
	require.NoError(t, err)

	srv := services.NewBlogService(db)

	params := &types.CreateBlogPostInput{
		Title:    "Test Post",
		Slug:     "test-post",
		Language: "en",
		Body:     new("Some body"),
	}

	t.Run("creates a blog post with valid params", func(t *testing.T) {
		post, err := srv.CreateBlogPost(ctx, params)
		assert.NoError(t, err)
		require.NotNil(t, post)
		assert.Equal(t, "test-post", post.Slug)
		assert.Nil(t, post.PublishedAt)
	})

	t.Run("rejects a duplicate slug", func(t *testing.T) {
		_, err := srv.CreateBlogPost(ctx, params)
		verrs, ok := errors.AsType[validation.Errors](err)
		require.True(t, ok, "expected validation.Errors, got %v", err)
		assert.Error(t, verrs["slug"])
	})

	t.Run("rejects invalid params", func(t *testing.T) {
		p := *params
		p.Slug = "another-slug"
		p.Title = ""

		_, err := srv.CreateBlogPost(ctx, &p)
		verrs, ok := errors.AsType[validation.Errors](err)
		require.True(t, ok, "expected validation.Errors, got %v", err)
		assert.Error(t, verrs["title"])
	})
}

func TestBlogService_PublishBlogPost(t *testing.T) {
	ctx := t.Context()
	db, err := initDB(ctx)
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Exec(ctx, "truncate blog_posts cascade")
	require.NoError(t, err)

	srv := services.NewBlogService(db)

	t.Run("publishes a draft post with a body", func(t *testing.T) {
		post, err := mocks.BlogPost(db, ctx)
		require.NoError(t, err)

		actual, err := srv.PublishBlogPost(ctx, post.ID)
		assert.NoError(t, err)
		require.NotNil(t, actual)
		assert.NotNil(t, actual.PublishedAt)
	})

	t.Run("rejects a post without a body", func(t *testing.T) {
		post, err := mocks.BlogPost(db, ctx, func(p *queries.InsertBlogPostParams) {
			p.Body = nil
		})
		require.NoError(t, err)

		_, err = srv.PublishBlogPost(ctx, post.ID)
		verrs, ok := errors.AsType[validation.Errors](err)
		require.True(t, ok, "expected validation.Errors, got %v", err)
		assert.Error(t, verrs["body"])
	})

	t.Run("rejects a post that is already published", func(t *testing.T) {
		post, err := mocks.BlogPost(db, ctx)
		require.NoError(t, err)

		_, err = srv.PublishBlogPost(ctx, post.ID)
		require.NoError(t, err)

		_, err = srv.PublishBlogPost(ctx, post.ID)
		verrs, ok := errors.AsType[validation.Errors](err)
		require.True(t, ok, "expected validation.Errors, got %v", err)
		assert.Error(t, verrs["publishedAt"])
	})

	t.Run("returns sql.ErrNoRows for an unknown ID", func(t *testing.T) {
		_, err := srv.PublishBlogPost(ctx, uuid.Must(uuid.NewV7()))
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})
}

func TestBlogService_UnpublishBlogPost(t *testing.T) {
	ctx := t.Context()
	db, err := initDB(ctx)
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Exec(ctx, "truncate blog_posts cascade")
	require.NoError(t, err)

	srv := services.NewBlogService(db)

	t.Run("unpublishes a published post", func(t *testing.T) {
		post, err := mocks.BlogPost(db, ctx)
		require.NoError(t, err)

		_, err = srv.PublishBlogPost(ctx, post.ID)
		require.NoError(t, err)

		actual, err := srv.UnpublishBlogPost(ctx, post.ID)
		assert.NoError(t, err)
		require.NotNil(t, actual)
		assert.Nil(t, actual.PublishedAt)
	})

	t.Run("rejects a post that is already a draft", func(t *testing.T) {
		post, err := mocks.BlogPost(db, ctx)
		require.NoError(t, err)

		_, err = srv.UnpublishBlogPost(ctx, post.ID)
		verrs, ok := errors.AsType[validation.Errors](err)
		require.True(t, ok, "expected validation.Errors, got %v", err)
		assert.Error(t, verrs["publishedAt"])
	})

	t.Run("returns sql.ErrNoRows for an unknown ID", func(t *testing.T) {
		_, err := srv.UnpublishBlogPost(ctx, uuid.Must(uuid.NewV7()))
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})
}
