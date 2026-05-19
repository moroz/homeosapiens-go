package jobs

import "github.com/google/uuid"

type OrderEmailType int

const (
	OrderEmailTypeOrderConfirmation OrderEmailType = iota
	OrderEmailTypePaymentConfirmation
)

type UserEmailType int

const (
	UserEmailTypeEmailVerification UserEmailType = iota
	UserEmailTypePasswordReset
)

type SendOrderEmailArgs struct {
	OrderID   uuid.UUID      `json:"order_id"`
	EmailType OrderEmailType `json:"email_type"`
}

func (SendOrderEmailArgs) Kind() string {
	return "SendOrderEmail"
}

type SendUserEmailArgs struct {
	UserID    uuid.UUID     `json:"user_id"`
	EmailType UserEmailType `json:"email_type"`
}

func (SendUserEmailArgs) Kind() string {
	return "SendUserEmail"
}

type VacuumUserTokensArgs struct{}

func (VacuumUserTokensArgs) Kind() string {
	return "VacuumUserTokens"
}
