-- +goose Up
alter table users add column preferred_timezone_encrypted bytea;

-- +goose Down
alter table users drop column preferred_timezone_encrypted;
