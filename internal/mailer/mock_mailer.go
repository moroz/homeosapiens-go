package mailer

import (
	"context"

	"github.com/wneessen/go-mail"
)

type mockMailer struct {
	messages []*mail.Msg
}

func MockMailer() Mailer {
	return &mockMailer{}
}

func (m *mockMailer) Send(_ context.Context, msg *mail.Msg) error {
	m.messages = append(m.messages, msg)
	return nil
}
