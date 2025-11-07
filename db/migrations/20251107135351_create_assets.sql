-- +goose Up
-- +goose StatementBegin
create table assets (
    id uuid primary key default uuidv7(),
    object_key text not null,
    original_filename text,
    inserted_at timestamp(0) not null default (now() at time zone 'utc'),
    updated_at timestamp(0) not null default (now() at time zone 'utc')
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table assets;
-- +goose StatementEnd
