-- +goose Up
-- +goose StatementBegin
delete from cart_line_items where event_id = '00000000-0000-0000-0000-000000000000' or cart_id = '00000000-0000-0000-0000-000000000000';
alter table cart_line_items add constraint cart_line_items_non_zero_foreign_keys
check (cart_id != '00000000-0000-0000-0000-000000000000' and event_id != '00000000-0000-0000-0000-000000000000');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table cart_line_items drop constraint cart_line_items_non_zero_foreign_keys;
-- +goose StatementEnd
