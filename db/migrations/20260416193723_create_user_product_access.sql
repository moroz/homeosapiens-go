-- +goose Up
-- +goose StatementBegin
create table user_product_access (
  id uuid not null default uuidv7(),
  user_id uuid not null references users (id) on delete cascade,
  product_id uuid not null references products (id) on delete cascade,
  order_id uuid references orders (id),
  inserted_at timestamp not null default now(),
  unique (user_id, product_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table user_product_access;
-- +goose StatementEnd
