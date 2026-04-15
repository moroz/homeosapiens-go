-- +goose Up
create type product_type as enum ('event', 'book');

create table products (
  id uuid not null primary key default uuidv7(),
  product_type product_type not null,
  title_pl text not null,
  title_en text not null,
  base_price_amount decimal not null default 0,
  base_price_currency char(3) not null default 'PLN',
  inserted_at timestamp not null default now(),
  updated_at timestamp not null default now()
);

truncate events cascade;
alter table events add column product_id uuid not null references products (id) on delete cascade, drop column base_price_amount, drop column base_price_currency;

alter table cart_line_items drop column event_id,
add column product_id uuid not null references products (id) on delete cascade;

alter table order_line_items drop column event_id,
add column product_id uuid not null references products (id) on delete cascade;

-- +goose Down
alter table order_line_items drop column product_id,
add column event_id uuid not null references events (id) on delete cascade;
alter table cart_line_items drop column product_id,
add column event_id uuid not null references events (id) on delete cascade;
alter table events drop column product_id;
drop table products;
drop type product_type;
