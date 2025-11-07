package helpers

import "github.com/nicksnyder/go-i18n/v2/i18n"

func TranslateCountry(localizer *i18n.Localizer, countryCode string) string {
	return localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "countries." + countryCode,
			Other: countryCode,
		},
	})
}

func TranslateSalutation(localizer *i18n.Localizer, salutation *string) string {
	if salutation == nil {
		return ""
	}

	return localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{
		ID:    *salutation,
		Other: *salutation,
	}}) + " "
}
