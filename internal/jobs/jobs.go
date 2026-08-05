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

type SendEventRegistrationEmailArgs struct {
	UserID  uuid.UUID `json:"user_id"`
	EventID uuid.UUID `json:"event_id"`
}

func (SendEventRegistrationEmailArgs) Kind() string {
	return "SendEventRegistrationEmail"
}

// ReminderLead identifies which of the two pre-event reminders a job is for.
type ReminderLead string

const (
	ReminderLead24h ReminderLead = "24h"
	ReminderLead1h  ReminderLead = "1h"
)

// EnqueueEventRemindersArgs runs periodically, claims the events whose reminder
// is due and fans out one SendEventReminderEmail job per registration.
type EnqueueEventRemindersArgs struct{}

func (EnqueueEventRemindersArgs) Kind() string {
	return "EnqueueEventReminders"
}

type SendEventReminderEmailArgs struct {
	UserID  uuid.UUID    `json:"user_id"`
	EventID uuid.UUID    `json:"event_id"`
	Lead    ReminderLead `json:"lead"`
}

func (SendEventReminderEmailArgs) Kind() string {
	return "SendEventReminderEmail"
}

type VacuumUserTokensArgs struct{}

func (VacuumUserTokensArgs) Kind() string {
	return "VacuumUserTokens"
}
