package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/moroz/securecookie"
)

type preferencesController struct {
	sessionStore securecookie.Store
}

func PreferencesController(sessionStore securecookie.Store) *preferencesController {
	return &preferencesController{
		sessionStore,
	}
}

func (c *preferencesController) SaveTimezone(w http.ResponseWriter, r *http.Request) {
	tzParam := r.URL.Query().Get("tz")
	if _, err := time.LoadLocation(tzParam); err != nil || tzParam == "" {
		http.Error(w, "Invalid timezone", 400)
	}
	session := r.Context().Value(config.SessionContextName).(types.SessionData)
	session["tz"] = tzParam
	if err := SaveSession(w, c.sessionStore, session); err != nil {
		log.Printf("Error serializing session cookie: %s", err)
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
