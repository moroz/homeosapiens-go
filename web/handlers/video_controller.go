package handlers

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"sync"

	_ "golang.org/x/image/webp"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/videos"
	"github.com/moroz/homeosapiens-go/web/helpers"
	"golang.org/x/image/draw"
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

// photoCache caches resized+encoded portrait data URIs keyed by CDN URL.
var photoCache sync.Map

// displayW/H are the SVG display dimensions at 2× (retina).
const portraitDisplayW = 480
const portraitDisplayH = 656

func fetchPortraitDataURI(ctx context.Context, url string) (string, error) {
	if cached, ok := photoCache.Load(url); ok {
		return cached.(string), nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	src, _, err := image.Decode(bytes.NewReader(body))
	if err != nil {
		// Unsupported format — embed raw without resizing
		ct := resp.Header.Get("Content-Type")
		if ct == "" {
			ct = "image/jpeg"
		}
		dataURI := "data:" + ct + ";base64," + base64.StdEncoding.EncodeToString(body)
		photoCache.Store(url, dataURI)
		return dataURI, nil
	}

	dst := image.NewRGBA(image.Rect(0, 0, portraitDisplayW, portraitDisplayH))
	draw.BiLinear.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Src, nil)

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, dst, &jpeg.Options{Quality: 82}); err != nil {
		return "", err
	}

	dataURI := "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
	photoCache.Store(url, dataURI)
	return dataURI, nil
}

func (cc *videoController) Thumbnail(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.ErrNotFound
	}
	locale := c.Param("locale")

	data, err := cc.videoService.GetVideoThumbnailData(c.Request().Context(), id)
	if err != nil {
		return err
	}

	title := data.Video.TitleEn
	if locale == "pl" {
		title = data.Video.TitlePl
	}

	var host string
	if data.HostGivenName != nil && data.HostFamilyName != nil {
		host = *data.HostGivenName + " " + *data.HostFamilyName
	}

	var ppURL *string
	if data.HostProfilePictureUrl != nil {
		cdnURL := fmt.Sprintf("%s/%s", config.AssetCdnBaseUrl, *data.HostProfilePictureUrl)
		dataURI, err := fetchPortraitDataURI(c.Request().Context(), cdnURL)
		if err == nil {
			ppURL = &dataURI
		}
	}

	c.Response().Header().Set("Content-Type", "image/svg+xml")
	c.Response().Header().Set("Cache-Control", "public, max-age=3600")
	return videos.Thumbnail(title, host, locale, data.Video.RecordedOn, ppURL).Render(c.Response())
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

func (cc *videoController) Youtube(c *echo.Context) error {
	ctx := helpers.GetRequestContext(c)

	data, err := queries.New(cc.db).ListYoutubeVideos(c.Request().Context())
	if err != nil {
		return err
	}

	return videos.Youtube(ctx, data).Render(c.Response())
}
