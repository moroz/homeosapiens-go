-- +goose Up
-- +goose StatementBegin
create unique index on cart_line_items (event_id, cart_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index cart_line_items_event_id_cart_id_idx;
-- +goose StatementEnd
