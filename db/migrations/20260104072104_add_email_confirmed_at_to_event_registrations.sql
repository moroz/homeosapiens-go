-- +goose Up
-- +goose StatementBegin
alter table event_registrations add column email_confirmed_at timestamp;
alter table users add column email_confirmed_at timestamp;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table event_registrations drop column email_confirmed_at;
alter table users drop column email_confirmed_at;
-- +goose StatementEnd
