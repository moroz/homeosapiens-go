-- +goose Up
create type product_type as enum ('event', 'book');

create table products (
  id uuid not null primary key default uuidv7(),
  product_type product_type not null,
  title_pl text not null,
  title_en text not null,
  base_price_amount decimal not null default 0,
  base_price_currency text not null default 'PLN',
  inserted_at timestamp not null default now(),
  updated_at timestamp not null default now()
);

truncate events cascade;
alter table events add product_id uuid references products (id) on delete cascade, drop base_price_amount, drop column base_price_currency;

alter table cart_line_items drop column event_id,
add product_id uuid not null references products (id) on delete cascade;

create unique index on cart_line_items (cart_id, product_id);

alter table order_line_items drop column event_id,
add product_id uuid not null references products (id) on delete cascade;
alter table order_line_items rename event_title to product_title;
alter table order_line_items rename event_price_amount to product_price_amount;
alter table order_line_items rename event_price_currency to product_price_currency;

truncate event_prices;
alter table event_prices drop column event_id, add column product_id uuid not null references products (id) on delete cascade;
create unique index on event_prices (discount_code, product_id);
alter table event_prices rename to product_prices;

-- +goose Down
truncate product_prices;
alter table product_prices rename to event_prices;
alter table event_prices drop column product_id, add column event_id uuid not null references events (id) on delete cascade;
create unique index on event_prices (discount_code, event_id);

alter table order_line_items drop column product_id,
add column event_id uuid not null references events (id) on delete cascade;
alter table order_line_items rename product_title to event_title;
alter table order_line_items rename product_price_amount to event_price_amount;
alter table order_line_items rename product_price_currency to event_price_currency;

alter table cart_line_items drop column product_id,
add column event_id uuid not null references events (id) on delete cascade;
create unique index on cart_line_items (cart_id, event_id);

alter table events drop column product_id, add column base_price_amount decimal(20,8), add column base_price_currency text;
drop table products;
drop type product_type;
