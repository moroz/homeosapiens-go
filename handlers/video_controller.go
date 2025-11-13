package handlers

import (
	"net/http"

	"github.com/moroz/homeosapiens-go/tmpl/videos"
)

type videoController struct{}

func VideoController() *videoController {
	return &videoController{}
}

func (c *videoController) Index(w http.ResponseWriter, r *http.Request) {
	if err := videos.Index(r.Context()).Render(w); err != nil {
		handleRenderingError(w, err)
	}
}
