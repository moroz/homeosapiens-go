-- +goose Up
-- +goose StatementBegin
alter table orders add column email_encrypted bytea not null;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table orders drop column email_encrypted;
-- +goose StatementEnd
