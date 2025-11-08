package helpers

import (
	"time"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

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

	translated, err := localizer.Localize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{
		ID:    *salutation,
		Other: *salutation,
	}})

	if err != nil {
		return *salutation + " "
	}

	return translated + " "
}

func FormatDate(date time.Time, locale string) string {
	switch locale {
	case "pl":
		return date.Format("02.01.2006")
	default:
		return date.Format("02/01/2006")
	}
}

func FormatDateRange(start, end time.Time, tz string, locale string) string {
	if tz == "" {
		tz = "Europe/Warsaw"
	}
	timezone, _ := time.LoadLocation(tz)

	start = start.In(timezone)
	end = end.In(timezone)
	d1, m1, y1 := start.In(timezone).Date()
	d2, m2, y2 := end.In(timezone).Date()

	if d1 == d2 && m1 == m2 && y1 == y2 {
		return FormatDate(start, locale)
	}

	switch locale {
	case "pl":
		return start.Format("02.01–") + FormatDate(end, locale)
	default:
		return start.Format("02/01–") + FormatDate(end, locale)
	}
}
