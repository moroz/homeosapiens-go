package mailers

import (
	"context"
	"fmt"

	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/tmpl/email"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type UserMailer interface {
	SendEmailVerification(context.Context, *queries.User) error
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

func (m userMailer) SendEmailVerification(ctx context.Context, user *queries.User) error {
	l := i18n.NewLocalizer(m.bundle, string(user.PreferredLocale))

	subject, err := l.Localize(&i18n.LocalizeConfig{
		MessageID:    "emails.email_verification.subject",
		TemplateData: user,
	})
	if err != nil {
		return fmt.Errorf("SendEmailVerification: %w", err)
	}

	props := email.EmailVerificationEmailProps{
		LayoutProps: &email.LayoutProps{
			Title:     subject,
			Language:  string(user.PreferredLocale),
			Localizer: l,
		},
	}
	panic("implement me")
}
