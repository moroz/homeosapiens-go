-- +goose Up
-- +goose StatementBegin
create table user_tokens (
    id uuid primary key default uuidv7(),
    user_id uuid not null references users (id),
    context text not null,
    token bytea not null unique,
    inserted_at timestamp(0) not null default (now() at time zone 'utc'),
    valid_until timestamp(0) not null
);

create index on user_tokens (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table user_tokens;
-- +goose StatementEnd
