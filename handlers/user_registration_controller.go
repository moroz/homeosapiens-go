package handlers

import (
	"net/http"

	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	userregistrations "github.com/moroz/homeosapiens-go/tmpl/user_registrations"
)

type userRegistrationController struct {
	*services.UserService
}

func UserRegistrationController(db queries.DBTX) *userRegistrationController {
	return &userRegistrationController{
		services.NewUserService(db),
	}
}

func (c *userRegistrationController) New(w http.ResponseWriter, r *http.Request) {
	if err := userregistrations.New(r.Context()).Render(w); err != nil {
		handleRenderingError(w, err)
	}
}
