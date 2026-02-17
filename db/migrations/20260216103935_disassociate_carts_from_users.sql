-- +goose Up
-- +goose StatementBegin
alter table cart_line_items drop constraint cart_line_items_cart_id_fkey;
drop table carts;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
create table carts (
    id uuid not null primary key default uuidv7(),
    owner_id uuid references users (id) on delete cascade,
    inserted_at timestamp not null default (now() at time zone 'utc'),
    updated_at timestamp not null default (now() at time zone 'utc')
);

alter table cart_line_items add constraint cart_line_items_cart_id_fkey foreign key (cart_id) references carts (id);
create unique index on carts (owner_id) where owner_id is not null;
-- +goose StatementEnd
