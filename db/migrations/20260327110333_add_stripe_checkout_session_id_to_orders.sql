-- +goose Up
alter table orders add column stripe_checkout_session_id text;

-- +goose Down
alter table orders drop column stripe_checkout_session_id;
