-- +goose Up
-- +goose StatementBegin
alter table users drop column email,
    drop column given_name,
    drop column family_name;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table users add column email citext unique,
                  add column given_name text,
                  add column family_name text;
-- +goose StatementEnd
