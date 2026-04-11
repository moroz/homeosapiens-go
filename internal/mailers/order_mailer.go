package mailers

import (
	"context"
	"fmt"
	"log"

	"github.com/moroz/homeosapiens-go/tmpl/email"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type OrderMailer interface {
	SendOrderConfirmation(context.Context, *types.OrderDTO) error
	SendPaymentConfirmation(context.Context, *types.OrderDTO) error
}

type orderMailer struct {
	Mailer
	bundle *i18n.Bundle
}

func NewOrderMailer(m Mailer, bundle *i18n.Bundle) OrderMailer {
	return &orderMailer{
		Mailer: m,
		bundle: bundle,
	}
}

func orderEmailProps(l *i18n.Localizer, subject string, order *types.OrderDTO) *email.OrderEmailProps {
	return &email.OrderEmailProps{
		LayoutProps: &email.LayoutProps{
			Title:     subject,
			Localizer: l,
			Language:  string(order.PreferredLocale),
		},
		Order: order,
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

	msg := NewMessage()
	props := orderEmailProps(l, subject, order)

	msg.Subject(subject)
	msg.SetBodyHTMLTemplate(email.OrderConfirmationTemplate, props)
	msg.To(order.EmailRecipient())

	err = m.Mailer.Send(ctx, msg)
	if err != nil {
		log.Printf("SendOrderConfirmation for order %v: %s", order.OrderNumber, err)
	}
	return err
}

func (m *orderMailer) SendPaymentConfirmation(ctx context.Context, order *types.OrderDTO) error {
	l := i18n.NewLocalizer(m.bundle, string(order.PreferredLocale))

	subject, err := l.Localize(&i18n.LocalizeConfig{
		MessageID:    "emails.payment_confirmation.subject",
		TemplateData: order.Order,
	})
	if err != nil {
		return fmt.Errorf("SendPaymentConfirmation: %w", err)
	}

	msg := NewMessage()
	props := orderEmailProps(l, subject, order)

	msg.Subject(subject)
	msg.SetBodyHTMLTemplate(email.PaymentConfirmationTemplate, props)
	msg.To(order.EmailRecipient())

	err = m.Mailer.Send(ctx, msg)
	if err != nil {
		log.Printf("SendPaymentConfirmation for order %v: %s", order.OrderNumber, err)
	}
	return err
}
