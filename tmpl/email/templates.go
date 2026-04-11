package email

import (
	"embed"
	"html/template"
	"log"

	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/tmpl/helpers"
	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/shopspring/decimal"

	"github.com/moroz/homeosapiens-go/types"
)

//go:embed *.html.tmpl
var templateFS embed.FS

type LayoutProps struct {
	Title     string
	Language  string
	Localizer *goi18n.Localizer
}

func (p *LayoutProps) LogoURL() string {
	return config.PublicUrl + "/assets/logo.png"
}

func (p *LayoutProps) T(messageID string) string {
	if p.Localizer == nil {
		return messageID
	}
	msg, err := p.Localizer.Localize(&goi18n.LocalizeConfig{MessageID: messageID})
	if err != nil {
		log.Print(err)
	}
	return msg
}

func (p *LayoutProps) Translate(messageID string, data any) template.HTML {
	if p.Localizer == nil {
		return template.HTML(messageID)
	}
	msg, err := p.Localizer.Localize(&goi18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: data,
	})
	if err != nil {
		log.Print(err)
	}
	return template.HTML(msg)
}

func (p *LayoutProps) FormatPrice(amount decimal.Decimal, currency string) string {
	return helpers.FormatPrice(amount, currency, p.Language)
}

type OrderEmailProps struct {
	*LayoutProps
	Order *types.OrderDTO
}

type UserEmailVerificationEmailProps struct {
	*LayoutProps
	UserToken *types.UserTokenDTO
}

var OrderConfirmationTemplate = template.Must(template.ParseFS(templateFS, "layout.html.tmpl", "_header.html.tmpl", "_footer.html.tmpl", "_order_summary.html.tmpl", "order_confirmation.html.tmpl"))

var PaymentConfirmationTemplate = template.Must(template.ParseFS(templateFS, "layout.html.tmpl", "_header.html.tmpl", "_footer.html.tmpl", "_order_summary.html.tmpl", "payment_confirmation.html.tmpl"))

var UserEmailVerificationTemplate = template.Must(template.ParseFS(templateFS, "layout.html.tmpl", "_header.html.tmpl", "_footer.html.tmpl", "user_email_verification.html.tmpl"))
