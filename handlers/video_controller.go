package handlers

import (
	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/videos"
	"github.com/moroz/homeosapiens-go/types"
)

type videoController struct {
	db           queries.DBTX
	videoService *services.VideoService
}

func VideoController(db queries.DBTX) *videoController {
	return &videoController{
		db:           db,
		videoService: services.NewVideoService(db),
	}
}

func (c *videoController) Index(r *echo.Context) error {
	ctx := r.Get("context").(*types.CustomContext)
	user := ctx.User

	data, err := c.videoService.ListVideosWithSources(r.Request().Context(), user)
	if err != nil {
		return err
	}

	return videos.Index(ctx, data).Render(r.Response())
}
