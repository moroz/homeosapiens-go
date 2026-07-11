package handlers

import (
	"fmt"
	"net/http"

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

func (cc *videoController) Thumbnail(c *echo.Context) error {
	id := c.Param("id")
	locale := c.Param("locale")

	// TODO: fetch video metadata and generate real SVG
	svg := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="320" height="180"><text>%s %s</text></svg>`, id, locale)
	c.Response().Header().Set("Cache-Control", "public, max-age=3600")
	return c.Blob(http.StatusOK, "image/svg+xml", []byte(svg))
}

func (cc *videoController) Index(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)

	data, err := cc.videoService.ListVideoGroupsForUser(c.Request().Context(), ctx.User.ID)
	if err != nil {
		return err
	}

	var groupSlug *string
	if param := c.Param("group_slug"); param != "" {
		groupSlug = &param
	}
	group, err := cc.videoService.GetVideoGroupDetails(c.Request().Context(), ctx.User.ID, groupSlug)
	return videos.Index(ctx, data, group).Render(c.Response())
}

func (cc *videoController) Show(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)

	video, err := cc.videoService.GetVideoForUser(c.Request().Context(), ctx.User.ID, c.Param("group_slug"), c.Param("video_slug"))
	if err != nil {
		return err
	}

	group, err := cc.videoService.GetVideoGroupDetails(c.Request().Context(), ctx.User.ID, new(c.Param("group_slug")))
	if err != nil {
		return err
	}

	return videos.Show(ctx, group, video).Render(c.Response())
}
