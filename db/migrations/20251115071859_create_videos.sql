-- +goose Up
-- +goose StatementBegin
create type video_provider as enum ('youtube', 'cloudfront');

create table videos (
    id uuid primary key default uuidv7(),
    provider video_provider not null,
    is_public boolean not null default false,
    title_en varchar(255) not null,
    title_pl varchar(255) not null,
    slug citext not null unique,
    object_key text not null,
    inserted_at timestamp(0) not null default (now() at time zone 'utc'),
    updated_at timestamp(0) not null default (now() at time zone 'utc')
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table videos;
-- +goose StatementEnd
