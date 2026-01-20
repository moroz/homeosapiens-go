package types

import (
	"net/url"
	"time"

	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type CustomContext struct {
	User        *queries.User
	Session     SessionData
	Localizer   *i18n.Localizer
	Language    string
	Timezone    *time.Location
	TimezoneSet bool
	RequestUrl  *url.URL
}
