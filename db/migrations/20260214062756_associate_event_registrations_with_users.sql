-- +goose Up
-- +goose StatementBegin
drop table event_registrations;

create table event_registrations (
    id          uuid primary key      default uuidv7(),
    event_id    uuid         not null references events (id),
    user_id     uuid         not null references users (id),
    inserted_at timestamp(0) not null default (now() at time zone 'utc')
);

create unique index on event_registrations (event_id, user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table event_registrations;

create table event_registrations (
    id uuid primary key default uuidv7(),
    event_id uuid not null references events (id),
    user_id uuid not null references users (id),
    attending_in_person boolean not null default false,
    is_host boolean not null default false,
    inserted_at timestamp(0) not null default (now() at time zone 'utc'),
    updated_at timestamp(0) not null default (now() at time zone 'utc')
);

create index on event_registrations (user_id);
-- +goose StatementEnd
