-- +goose Up
-- +goose StatementBegin
alter table event_registrations add column licence_number_encrypted bytea;
alter table users add column licence_number_encrypted bytea;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table event_registrations drop column licence_number_encrypted;
alter table users drop column licence_number_encrypted;
-- +goose StatementEnd
