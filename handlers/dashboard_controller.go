package handlers

import (
	"net/http"

	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/tmpl/dashboard"
)

type dashboardController struct {
	db queries.DBTX
}

func DashboardController(db queries.DBTX) *dashboardController {
	return &dashboardController{db}
}

func (c *dashboardController) Index(w http.ResponseWriter, r *http.Request) {
	if err := dashboard.Index(r.Context()).Render(w); err != nil {
		handleRenderingError(w, err)
	}
}
