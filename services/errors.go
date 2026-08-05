package services

import (
	"errors"
	"fmt"
	"time"
)

var ErrTokenInvalid = errors.New("VerifyEmailAddress: token invalid or expired")
var ErrRateLimited = errors.New("rate limited")
var ErrEmailAlreadyVerified = errors.New("email address is already verified")

// ErrEventHasRegistrations reports an attempt to delete an event that people
// have already signed up for. Registrations are the record of who is owed a
// seat, so the event has to be unpublished instead.
var ErrEventHasRegistrations = errors.New("event has registrations")

type errRateLimited struct {
	limitedUntil time.Time
}

func (e errRateLimited) Error() string {
	return fmt.Sprintf("rate limited until %s", e.limitedUntil.Format(time.RFC3339))
}

func (e errRateLimited) Unwrap() error {
	return ErrRateLimited
}
