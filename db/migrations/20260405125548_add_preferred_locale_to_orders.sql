-- +goose Up
create type locale as enum ('pl', 'en');

alter table orders add column preferred_locale locale not null default 'en';
alter table users add column preferred_locale locale not null default 'en';

-- +goose Down
alter table orders drop column preferred_locale;
alter table users drop column preferred_locale;

drop type locale;
