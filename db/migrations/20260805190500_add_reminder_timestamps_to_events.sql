-- +goose Up
alter table events
add reminder_24h_sent_at timestamp(0),
add reminder_1h_sent_at timestamp(0);

-- +goose Down
alter table events
drop reminder_24h_sent_at,
drop reminder_1h_sent_at;
