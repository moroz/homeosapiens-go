package types

import (
	"net/url"
	"time"

	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

//go:generate stringer -type=FlashLevel -trimprefix=FlashLevel_
type FlashLevel int

const (
	FlashLevel_Info FlashLevel = iota
	FlashLevel_Success
	FlashLevel_Error
)

type FlashMessage struct {
	Level   FlashLevel
	Message string
}

type CustomContext struct {
	User        *queries.User
	Session     SessionData
	Localizer   *i18n.Localizer
	Language    string
	Timezone    *time.Location
	TimezoneSet bool
	RequestUrl  *url.URL
	Flash       []FlashMessage
}

func (c *CustomContext) IsPolish() bool {
	return c.Language == "pl"
}
