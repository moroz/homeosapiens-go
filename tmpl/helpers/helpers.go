package helpers

import (
	"fmt"
	"strings"
	"time"

	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/shopspring/decimal"
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

	translated, err := localizer.LocalizeMessage(&i18n.Message{
		ID:    *salutation,
		Other: *salutation,
	})

	if err != nil {
		return *salutation + " "
	}

	return translated + " "
}

func TranslateEventType(localizer *i18n.Localizer, eventType queries.EventType) string {
	translated, _ := localizer.LocalizeMessage(&i18n.Message{
		ID:    "common.events.event_type." + string(eventType),
		Other: string(eventType),
	})

	return translated
}

func FormatDate(date time.Time, locale string) string {
	switch locale {
	case "pl":
		return date.Format("02.01.2006")
	default:
		return date.Format("02/01/2006")
	}
}

func FormatDateTime(ts time.Time, locale string) string {
	switch locale {
	case "pl":
		return ts.Format("02.01.2006 15:04")
	default:
		return ts.Format("02 Jan, 2006, 03:04 PM")
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

var CurrencyFormats = map[string]map[string]string{
	"pl": {
		"PLN": "%v zł",
	},
	"en": {
		"PLN": "PLN %v",
	},
}

func FormatPrice(amount decimal.Decimal, currencyCode string, locale string) string {
	if m, ok := CurrencyFormats[locale]; ok {
		if format, ok := m[currencyCode]; ok {
			return fmt.Sprintf(format, amount.String())
		}
	}
	return fmt.Sprintf("%s %s", currencyCode, amount)
}

func FormatHostName(localizer *i18n.Localizer, host *queries.ListHostsForEventsRow) string {
	salutation := TranslateSalutation(localizer, host.Salutation)
	return fmt.Sprintf("%s%s %s", salutation, host.GivenName, host.FamilyName)
}

func FormatHosts(localizer *i18n.Localizer, hosts []*queries.ListHostsForEventsRow) string {
	names := make([]string, len(hosts))
	for i, host := range hosts {
		names[i] = FormatHostName(localizer, host)
	}
	return strings.Join(names, ", ")
}
