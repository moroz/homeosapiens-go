-- +goose Up
-- +goose StatementBegin
create table users (
    id uuid primary key default uuidv7(),
    email citext not null unique,
    salutation varchar(20),
    given_name varchar(255) not null,
    family_name varchar(255) not null,
    country char(2),
    profession varchar(255),
    organization varchar(255),
    company varchar(255),
    password_hash varchar(255),
    last_login_at timestamp(0),
    last_login_ip inet,
    inserted_at timestamp(0) not null default (now() at time zone 'utc'),
    updated_at timestamp(0) not null default (now() at time zone 'utc')
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users;
-- +goose StatementEnd
