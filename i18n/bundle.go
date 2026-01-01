package i18n

import (
	"embed"
	"encoding/json"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed en.json pl.json
var LocaleFS embed.FS

func InitBundle() (*i18n.Bundle, error) {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	for _, lang := range SupportedLocales {
		if _, err := bundle.LoadMessageFileFS(LocaleFS, lang+".json"); err != nil {
			return nil, err
		}
	}

	return bundle, nil
}
