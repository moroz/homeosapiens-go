-- +goose Up
-- Reminders are per attendee, not per event: somebody registering twelve hours
-- before the start would otherwise never get a reminder, because the event was
-- already stamped when the day-ahead scan ran.
alter table event_registrations
add reminder_24h_sent_at timestamp(0),
add reminder_1h_sent_at timestamp(0);

create index on event_registrations (event_id) where reminder_24h_sent_at is null;
create index on event_registrations (event_id) where reminder_1h_sent_at is null;

alter table events
drop reminder_24h_sent_at,
drop reminder_1h_sent_at;

-- +goose Down
alter table events
add reminder_24h_sent_at timestamp(0),
add reminder_1h_sent_at timestamp(0);

alter table event_registrations
drop reminder_24h_sent_at,
drop reminder_1h_sent_at;
