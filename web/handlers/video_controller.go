package handlers

import (
	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/videos"
	"github.com/moroz/homeosapiens-go/web/helpers"
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

func (cc *videoController) Index(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)

	data, err := cc.videoService.ListVideosWithSources(c.Request().Context(), ctx.User)
	if err != nil {
		return err
	}

	return videos.Index(ctx, data).Render(c.Response())
}
