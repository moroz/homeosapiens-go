package i18n

import (
	"log"

	"golang.org/x/text/language"
)

var SupportedLocales = []string{"en", "pl"}

var DefaultLocale = "en"

var matcher language.Matcher

func init() {
	var supported []language.Tag
	for _, lang := range SupportedLocales {
		tag, err := language.Parse(lang)
		if err != nil {
			log.Fatalf("Failed to parse supported language %s: %s", lang, err)
		}
		supported = append(supported, tag)
	}

	matcher = language.NewMatcher(supported)
}

func ResolveLocaleFromAcceptLanguageHeader(acceptLanguage string) (string, error) {
	tags, _, err := language.ParseAcceptLanguage(acceptLanguage)
	if err != nil {
		return "", err
	}

	tag, _, confidence := matcher.Match(tags...)
	if confidence == language.No {
		return "", nil
	}

	base, _ := tag.Base()
	return base.String(), nil
}

func ResolveLocale(sources ...string) string {
	for _, source := range sources {
		guess, _ := ResolveLocaleFromAcceptLanguageHeader(source)
		if guess != "" {
			return guess
		}
	}
	return DefaultLocale
}
