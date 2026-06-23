-- +goose Up
alter table users add preferred_timezone_locked boolean not null default false;

-- +goose Down
alter table users drop preferred_timezone_locked;
