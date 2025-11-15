package handlers

import (
	"log"
	"net/http"

	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/tmpl/videos"
)

type videoController struct {
	db queries.DBTX
}

func VideoController(db queries.DBTX) *videoController {
	return &videoController{db}
}

func (c *videoController) Index(w http.ResponseWriter, r *http.Request) {
	videoRows, err := queries.New(c.db).ListVideos(r.Context())
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), 500)
		return
	}

	if err := videos.Index(r.Context(), videoRows).Render(w); err != nil {
		handleRenderingError(w, err)
	}
}
