package mailers

import (
	"context"
	"fmt"
	"log"

	"github.com/moroz/homeosapiens-go/tmpl/email"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type UserMailer interface {
	SendUserEmailVerification(context.Context, *types.UserTokenDTO) error
	SendUserPasswordResetEmail(context.Context, *types.UserTokenDTO) error
}

type userMailer struct {
	Mailer
	bundle *i18n.Bundle
}

func NewUserMailer(m Mailer, bundle *i18n.Bundle) UserMailer {
	return &userMailer{
		Mailer: m,
		bundle: bundle,
	}
}

func (m userMailer) SendUserEmailVerification(ctx context.Context, data *types.UserTokenDTO) error {
	l := i18n.NewLocalizer(m.bundle, string(data.PreferredLocale))

	subject, err := l.Localize(&i18n.LocalizeConfig{
		MessageID:    "emails.user_email_verification.subject",
		TemplateData: data,
	})
	if err != nil {
		return fmt.Errorf("SendEmailVerification: %w", err)
	}

	props := email.UserEmailVerificationEmailProps{
		LayoutProps: &email.LayoutProps{
			Title:     subject,
			Language:  string(data.PreferredLocale),
			Localizer: l,
		},
		UserToken: data,
	}

	msg := NewMessage()
	msg.Subject(subject)
	msg.SetBodyHTMLTemplate(email.UserEmailVerificationTemplate, props)
	msg.To(data.EmailRecipient())

	err = m.Mailer.Send(ctx, msg)
	if err != nil {
		log.Printf("SendUserEmailVerification for user %v: %s", data.User.ID, err)
	}
	return err
}

func (m userMailer) SendUserPasswordResetEmail(ctx context.Context, data *types.UserTokenDTO) error {
	l := i18n.NewLocalizer(m.bundle, string(data.PreferredLocale))
	subject, err := l.Localize(&i18n.LocalizeConfig{
		MessageID:    "emails.password_reset.subject",
		TemplateData: data,
	})
	if err != nil {
		return fmt.Errorf("SendPasswordResetEmail: %w", err)
	}

	props := email.UserPasswordResetEmailProps{
		LayoutProps: &email.LayoutProps{
			Title:     subject,
			Language:  string(data.PreferredLocale),
			Localizer: l,
		},
		UserToken: data,
	}

	msg := NewMessage()
	msg.Subject(subject)
	msg.SetBodyHTMLTemplate(email.UserPasswordResetTemplate, props)
	msg.To(data.EmailRecipient())

	err = m.Mailer.Send(ctx, msg)
	if err != nil {
		log.Printf("SendPasswordResetEmail for user %v: %s", data.User.ID, err)
	}
	return err
}
