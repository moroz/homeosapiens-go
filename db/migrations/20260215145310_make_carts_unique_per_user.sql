-- +goose Up
-- +goose StatementBegin
create unique index on carts (owner_id) where owner_id is not null;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index carts_owner_id_idx;
-- +goose StatementEnd
