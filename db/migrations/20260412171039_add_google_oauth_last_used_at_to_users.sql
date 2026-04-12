-- +goose Up
alter table users add column google_oauth_last_used_at timestamp;

-- +goose Down
alter table users drop column google_oauth_last_used_at;
