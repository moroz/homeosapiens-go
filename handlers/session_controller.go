package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/moroz/homeosapiens-go/config"
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

func (c *sessionController) New(w http.ResponseWriter, r *http.Request) {
	if err := sessions.New(r.Context(), "", "").Render(w); err != nil {
		handleRenderingError(w, err)
	}
}

func (c *sessionController) Create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("Error parsing form: %s", err)
		http.Error(w, err.Error(), 400)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := c.UserService.AuthenticateUserByEmailPassword(r.Context(), email, password)
	if err != nil {
		l := r.Context().Value("localizer").(*i18n.Localizer)
		msg := l.MustLocalizeMessage(&i18n.Message{
			ID: "sessions.new.invalid_email_password_combination",
		})

		if err := sessions.New(r.Context(), email, msg).Render(w); err != nil {
			msg := fmt.Sprintf("Error rendering page: %s", err)
			log.Print(msg)
			http.Error(w, msg, 500)
		}
		return
	}

	token, err := c.UserTokenService.IssueAccessTokenForUser(r.Context(), user, 24*time.Hour)
	if err != nil {
		log.Printf("Error issuing access token: %s", err)
		http.Error(w, err.Error(), 500)
		return
	}

	session := r.Context().Value(config.SessionContextName).(types.SessionData)
	session["access_token"] = token.Token
	if err := SaveSession(w, c.sessionStore, session); err != nil {
		log.Printf("Error serializing session cookie: %s", err)
		http.Error(w, err.Error(), 500)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (c *sessionController) Delete(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(config.SessionContextName).(types.SessionData)
	token, ok := session["access_token"].([]byte)
	if ok && token != nil {
		if _, err := c.UserTokenService.DeleteToken(r.Context(), token); err != nil {
			log.Printf("Error deleting user token: %s", err)
		}
	}
	delete(session, "access_token")
	if err := SaveSession(w, c.sessionStore, session); err != nil {
		log.Printf("Error serializing session cookie: %s", err)
		http.Error(w, err.Error(), 500)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
