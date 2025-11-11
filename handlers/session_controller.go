package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/sessions"
)

type sessionController struct {
	userService *services.UserService
}

func SessionController(db queries.DBTX) *sessionController {
	return &sessionController{
		userService: services.NewUserService(db),
	}
}

func (c *sessionController) New(w http.ResponseWriter, r *http.Request) {
	if err := sessions.New(r.Context()).Render(w); err != nil {
		msg := fmt.Sprintf("Error rendering page: %s", err)
		log.Print(msg)
		http.Error(w, msg, 500)
	}
}
