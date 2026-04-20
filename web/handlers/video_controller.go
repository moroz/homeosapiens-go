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

	data, err := cc.videoService.ListVideoGroupsForUser(c.Request().Context(), ctx.User)
	if err != nil {
		return err
	}

	var groupSlug *string
	if param := c.Param("group_slug"); param != "" {
		groupSlug = &param
	}
	group, err := cc.videoService.GetVideoGroupDetails(c.Request().Context(), ctx.User, groupSlug)
	return videos.Index(ctx, data, group).Render(c.Response())
}

func (cc *videoController) Show(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)

	video, err := cc.videoService.GetVideoForUser(c.Request().Context(), ctx.User, c.Param("group_slug"), c.Param("video_slug"))
	if err != nil {
		return err
	}

	group, err := cc.videoService.GetVideoGroupDetails(c.Request().Context(), ctx.User, new(c.Param("group_slug")))
	if err != nil {
		return err
	}

	return videos.Show(ctx, group, video).Render(c.Response())
}
