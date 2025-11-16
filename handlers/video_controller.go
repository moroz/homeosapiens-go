package handlers

import (
	"log"
	"net/http"

	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/videos"
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

func (c *videoController) Index(w http.ResponseWriter, r *http.Request) {
	data, err := c.videoService.ListVideosWithSources(r.Context())
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), 500)
		return
	}

	if err := videos.Index(r.Context(), data).Render(w); err != nil {
		handleRenderingError(w, err)
	}
}
