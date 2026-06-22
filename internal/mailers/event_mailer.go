package mailers

import (
	"bytes"
	"context"
	"fmt"
	"log"

	gomail "github.com/wneessen/go-mail"

	"github.com/moroz/homeosapiens-go/tmpl/email"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type EventMailer interface {
	SendEventRegistrationConfirmation(context.Context, *types.EventRegistrationEmailDTO) error
}

type eventMailer struct {
	Mailer
	bundle *i18n.Bundle
}

func NewEventMailer(m Mailer, bundle *i18n.Bundle) EventMailer {
	return &eventMailer{Mailer: m, bundle: bundle}
}

func (m *eventMailer) SendEventRegistrationConfirmation(ctx context.Context, data *types.EventRegistrationEmailDTO) error {
	lang := string(data.User.PreferredLocale)
	l := i18n.NewLocalizer(m.bundle, lang)

	subject, err := l.Localize(&i18n.LocalizeConfig{
		MessageID:    "emails.event_registration_confirmation.subject",
		TemplateData: data,
	})
	if err != nil {
		return fmt.Errorf("SendEventRegistrationConfirmation: %w", err)
	}

	props := &email.EventRegistrationEmailProps{
		LayoutProps: &email.LayoutProps{
			Title:     subject,
			Language:  lang,
			Localizer: l,
		},
		Data: data,
	}

	msg := NewMessage()
	msg.Subject(subject)
	msg.To(data.EmailRecipient())
	msg.SetBodyHTMLTemplate(email.EventRegistrationConfirmationTemplate, props)

	icsData := data.ICS()
	if err := msg.AttachReader("invite.ics", bytes.NewReader(icsData),
		gomail.WithFileContentType("text/calendar"),
		gomail.WithFileEncoding(gomail.EncodingB64),
	); err != nil {
		return fmt.Errorf("SendEventRegistrationConfirmation: failed to attach ICS: %w", err)
	}

	err = m.Mailer.Send(ctx, msg)
	if err != nil {
		log.Printf("SendEventRegistrationConfirmation for user %v: %s", data.User.ID, err)
	}
	return err
}
