-- +goose Up
create unique index on orders (stripe_checkout_session_id);

-- +goose Down
drop index orders_stripe_checkout_session_id_idx;
