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

type OrderMailer interface {
	SendOrderConfirmation(context.Context, *types.OrderDTO) error
}

type orderMailer struct {
	mailer.Mailer
	bundle *i18n.Bundle
}

func NewOrderMailer(m mailer.Mailer, bundle *i18n.Bundle) OrderMailer {
	return &orderMailer{
		Mailer: m,
		bundle: bundle,
	}
}

func (m *orderMailer) SendOrderConfirmation(ctx context.Context, order *types.OrderDTO) error {
	l := i18n.NewLocalizer(m.bundle, string(order.PreferredLocale))

	subject, err := l.Localize(&i18n.LocalizeConfig{
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
			Localizer: l,
			Language:  string(order.PreferredLocale),
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
