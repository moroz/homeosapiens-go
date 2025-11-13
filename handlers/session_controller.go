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
	"github.com/moroz/securecookie"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func handleRenderingError(w http.ResponseWriter, err error) {
	msg := fmt.Sprintf("Error rendering page: %s", err)
	log.Print(msg)
	http.Error(w, msg, 500)
}

type sessionController struct {
	userService      *services.UserService
	userTokenService *services.UserTokenService
	sessionStore     securecookie.Store
}

func SessionController(db queries.DBTX, sessionStore securecookie.Store) *sessionController {
	return &sessionController{
		userService:      services.NewUserService(db),
		userTokenService: services.NewUserTokenService(db),
		sessionStore:     sessionStore,
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

	user, err := c.userService.AuthenticateUserByEmailPassword(r.Context(), email, password)
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

	token, err := c.userTokenService.IssueAccessTokenForUser(r.Context(), user, 24*time.Hour)
	if err != nil {
		log.Printf("Error parsing form: %s", err)
		http.Error(w, err.Error(), 500)
		return
	}

	sessionData := r.Context().Value(config.SessionContextName).(*SessionData)
	sessionData.AccessToken = token.Token
	if err := SaveSession(w, c.sessionStore, sessionData); err != nil {
		log.Printf("Error serializing session cookie: %s", err)
		http.Error(w, err.Error(), 500)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
