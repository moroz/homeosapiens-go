package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/sessions"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/moroz/securecookie"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func handleError(w http.ResponseWriter, err error, code int) {
	msg := fmt.Sprintf("Error: %s", err)
	log.Print(msg)
	http.Error(w, msg, code)
}

func handleRenderingError(w http.ResponseWriter, err error) {
	msg := fmt.Sprintf("Error rendering page: %s", err)
	log.Print(msg)
	http.Error(w, msg, 500)
}

type sessionController struct {
	*services.UserService
	*services.UserTokenService
	sessionStore securecookie.Store
}

func SessionController(db queries.DBTX, sessionStore securecookie.Store) *sessionController {
	return &sessionController{
		services.NewUserService(db),
		services.NewUserTokenService(db),
		sessionStore,
	}
}

func (c *sessionController) New(r *echo.Context) error {
	ctx := r.Get("context").(*types.CustomContext)
	return sessions.New(ctx, "", "").Render(r.Response())
}

func (c *sessionController) Create(r *echo.Context) error {
	ctx := r.Get("context").(*types.CustomContext)

	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := c.UserService.AuthenticateUserByEmailPassword(r.Request().Context(), email, password)
	if err != nil {
		l := ctx.Localizer
		msg := l.MustLocalizeMessage(&i18n.Message{
			ID: "sessions.new.invalid_email_password_combination",
		})

		return sessions.New(ctx, email, msg).Render(r.Response())
	}

	token, err := c.UserTokenService.IssueAccessTokenForUser(r.Request().Context(), user, 24*time.Hour)
	if err != nil {
		return err
	}

	session := ctx.Session
	session["access_token"] = token.Token
	if err := SaveSession(r.Response(), c.sessionStore, session); err != nil {
		return err
	}

	return r.Redirect(http.StatusFound, "/")
}

func (c *sessionController) Delete(r *echo.Context) error {
	ctx := r.Get("context").(*types.CustomContext)
	session := ctx.Session
	token, ok := session["access_token"].([]byte)
	if ok && token != nil {
		if _, err := c.UserTokenService.DeleteToken(r.Request().Context(), token); err != nil {
			log.Printf("Error deleting user token: %s", err)
		}
	}
	delete(session, "access_token")
	if err := SaveSession(r.Response(), c.sessionStore, session); err != nil {
		return err
	}

	return r.Redirect(http.StatusFound, "/")
}
