-- +goose Up
-- +goose StatementBegin
create table hosts (
    id uuid primary key default uuidv7(),
    salutation varchar(255),
    given_name varchar(255) not null,
    family_name varchar(255) not null,
    profile_picture_id uuid references assets (id) ON DELETE SET NULL,
    inserted_at timestamp(0) not null default (now() at time zone 'utc'),
    updated_at timestamp(0) not null default (now() at time zone 'utc')
);

create table events_hosts (
    id uuid primary key default uuidv7(),
    event_id uuid not null references events (id),
    host_id uuid not null references hosts (id),
    priority int not null,
    inserted_at timestamp(0) not null default (now() at time zone 'utc'),
    updated_at timestamp(0) not null default (now() at time zone 'utc')
);

create unique index on events_hosts (event_id, priority);
create unique index on events_hosts (event_id, host_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table events_hosts;
drop table hosts;
-- +goose StatementEnd
