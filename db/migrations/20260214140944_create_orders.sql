-- +goose Up
-- +goose StatementBegin
create table carts (
    id uuid not null primary key default uuidv7(),
    owner_id uuid references users (id) on delete cascade,
    inserted_at timestamp not null default (now() at time zone 'utc'),
    updated_at timestamp not null default (now() at time zone 'utc')
);

create table cart_line_items (
    id uuid not null primary key default uuidv7(),
    cart_id uuid not null references carts (id) on delete cascade,
    event_id uuid not null references events (id) on delete cascade,
    quantity integer not null default 1,
    inserted_at timestamp not null default (now() at time zone 'utc'),
    updated_at timestamp not null default (now() at time zone 'utc')
);

create table orders (
    id uuid not null primary key default uuidv7(),
    user_id uuid not null references users (id),
    paid_at timestamp,
    cancelled_at timestamp,
    discount_code text,
    grand_total decimal not null,
    currency char(3) not null default 'PLN',
    inserted_at timestamp not null default (now() at time zone 'utc'),
    updated_at timestamp not null default (now() at time zone 'utc')
);

create table order_line_items (
    id uuid not null primary key default uuidv7(),
    order_id uuid not null references orders (id) on delete cascade,
    event_id uuid not null references events (id) on delete set null,
    event_title text not null,
    event_price_amount decimal not null,
    event_price_currency char(3) not null default 'PLN',
    quantity integer not null default 1,
    inserted_at timestamp not null default (now() at time zone 'utc'),
    updated_at timestamp not null default (now() at time zone 'utc')
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table order_line_items, orders, cart_line_items, carts;
-- +goose StatementEnd
