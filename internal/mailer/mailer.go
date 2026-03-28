package mailer

import (
	"context"

	"github.com/wneessen/go-mail"
)

type Mailer interface {
	Send(context.Context, *mail.Msg) error
}

type SMTPMailer struct {
	client *mail.Client
}

func NewSMTPMailer(host string, port int, username, password string) (Mailer, error) {
	client, err := mail.NewClient(
		host,
		mail.WithPort(port),
		mail.WithSSLPort(false),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(username),
		mail.WithPassword(password),
	)

	if err != nil {
		return nil, err
	}

	return &SMTPMailer{
		client: client,
	}, nil
}

func (m *SMTPMailer) Send(ctx context.Context, msg *mail.Msg) error {
	return m.client.DialAndSendWithContext(ctx, msg)
}
