-- +goose Up
-- +goose StatementBegin
create type event_type as enum ('seminar', 'webinar');

CREATE TABLE events (
    id uuid primary key default uuidv7(),
    title_en text not null,
    title_pl text not null,
    starts_at timestamp(0) not null,
    ends_at timestamp(0) not null,
    is_virtual boolean not null default false,
    description_en text not null,
    description_pl text,
    event_type event_type not null default 'seminar',
    base_price_amount decimal(20,8),
    base_price_currency text,
    inserted_at timestamp(0) not null default (now() at time zone 'utc'),
    updated_at timestamp(0) not null default (now() at time zone 'utc'),
    check (ends_at > starts_at),
    check ((base_price_amount is null) = (base_price_currency is null))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table events;
drop type event_type;
-- +goose StatementEnd
