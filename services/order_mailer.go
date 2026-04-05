package services

import (
	"context"
	"fmt"

	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/internal/mailer"
	"github.com/moroz/homeosapiens-go/tmpl/email"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/wneessen/go-mail"
)

type OrderMailer struct {
	mailer.Mailer
	Language  string
	Localizer *i18n.Localizer
}

func NewOrderMailer(m mailer.Mailer, lang string, l *i18n.Localizer) *OrderMailer {
	return &OrderMailer{
		Mailer:    m,
		Language:  lang,
		Localizer: l,
	}
}

func (m *OrderMailer) SendOrderConfirmation(ctx context.Context, order *types.OrderDTO) error {
	subject, err := m.Localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    "emails.order_confirmation.subject",
		TemplateData: order.Order,
	})
	if err != nil {
		return fmt.Errorf("SendOrderConfirmation: %w", err)
	}

	props := email.OrderConfirmationEmailProps{
		LayoutProps: &email.LayoutProps{
			Title:     subject,
			LogoURL:   config.PublicUrl + "/assets/logo.png",
			Localizer: m.Localizer,
			Language:  m.Language,
		},
		Order: order,
	}

	msg := mail.NewMsg(
		mail.WithEncoding(mail.EncodingQP),
		mail.WithCharset(mail.CharsetUTF8),
	)

	msg.Subject(subject)
	msg.SetBodyHTMLTemplate(email.OrderConfirmationTemplate, props)
	msg.From("Homeo sapiens <homeo.zoom@gmail.com>")
	msg.To(fmt.Sprintf("%s %s <%s>", order.BillingGivenName.String(), order.BillingFamilyName.String(), order.Email.String()))

	return m.Mailer.Send(ctx, msg)
}
