package email

import (
	"embed"
	"html/template"

	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"

	"github.com/moroz/homeosapiens-go/types"
)

//go:embed *.html.tmpl
var templateFS embed.FS

type LayoutProps struct {
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

type OrderConfirmationEmailProps struct {
	*LayoutProps
	Order *types.OrderDTO
}

var LayoutTemplate = template.Must(template.ParseFS(templateFS, "layout.html.tmpl", "_header.html.tmpl", "_footer.html.tmpl"))

var OrderConfirmationTemplate = template.Must(template.ParseFS(templateFS, "layout.html.tmpl", "_header.html.tmpl", "_footer.html.tmpl", "order_confirmation.html.tmpl"))
