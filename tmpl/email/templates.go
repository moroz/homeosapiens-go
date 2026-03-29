package email

import (
	"embed"
	"html/template"

	"github.com/moroz/homeosapiens-go/tmpl/helpers"
	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/shopspring/decimal"

	"github.com/moroz/homeosapiens-go/types"
)

//go:embed *.html.tmpl
var templateFS embed.FS

type LayoutProps struct {
	Language  string
	LogoURL   string
	Localizer *goi18n.Localizer
}

func (p *LayoutProps) T(messageID string) string {
	if p.Localizer == nil {
		return messageID
	}
	msg, _ := p.Localizer.Localize(&goi18n.LocalizeConfig{MessageID: messageID})
	return msg
}

func (p *LayoutProps) FormatPrice(amount decimal.Decimal, currency string) string {
	return helpers.FormatPrice(amount, currency, p.Language)
}

type OrderConfirmationEmailProps struct {
	*LayoutProps
	Order *types.OrderDTO
}

var LayoutTemplate = template.Must(template.ParseFS(templateFS, "layout.html.tmpl", "_header.html.tmpl", "_footer.html.tmpl"))

var OrderConfirmationTemplate = template.Must(template.ParseFS(templateFS, "layout.html.tmpl", "_header.html.tmpl", "_footer.html.tmpl", "order_confirmation.html.tmpl"))
