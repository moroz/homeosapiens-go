package handlers

import (
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
