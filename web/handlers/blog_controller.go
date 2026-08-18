package handlers

import (
	"database/sql"
	"errors"

	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/tmpl/blog"
	"github.com/moroz/homeosapiens-go/web/helpers"
)

type blogController struct {
	db queries.DBTX
}

func BlogController(db queries.DBTX) *blogController {
	return &blogController{db}
}

func (cc *blogController) Index(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)

	posts, err := queries.New(cc.db).ListPublishedBlogPostsByLanguage(c.Request().Context(), queries.Locale(ctx.Language))
	if err != nil {
		return err
	}

	return wrapRender(blog.Index(ctx, posts), c.Response())
}

func (cc *blogController) Show(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)
	slug := c.Param("slug")

	post, err := queries.New(cc.db).GetBlogPostBySlug(c.Request().Context(), slug)
	if errors.Is(err, sql.ErrNoRows) {
		return echo.ErrNotFound
	}
	if err != nil {
		return err
	}

	// Unpublished posts are only reachable through the admin UI.
	if post.PublishedAt == nil {
		return echo.ErrNotFound
	}

	return wrapRender(blog.Show(ctx, post), c.Response())
}
