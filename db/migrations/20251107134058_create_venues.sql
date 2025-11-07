-- +goose Up
-- +goose StatementBegin
create table venues (
    id uuid primary key default uuidv7(),
    name_en varchar(255) not null,
    name_pl varchar(255),
    street varchar(255) not null,
    city_en varchar(255) not null,
    city_pl varchar(255),
    postal_code varchar(15),
    country_code varchar(2) not null,
    inserted_at timestamp not null default (now() at time zone 'utc'),
    updated_at timestamp not null default (now() at time zone 'utc')
);

alter table events add column venue_id uuid references venues (id) ON DELETE SET NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table events drop column venue_id;

drop table venues;
-- +goose StatementEnd
