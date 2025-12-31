-- +goose Up
-- +goose StatementBegin
truncate users cascade;

alter table users add column email_encrypted bytea,
    add column email_hash bytea unique,
    add column given_name_encrypted bytea,
    add column family_name_encrypted bytea;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table users drop column email_encrypted,
    drop column email_hash,
    drop column given_name_encrypted,
    drop column family_name_encrypted;
-- +goose StatementEnd
