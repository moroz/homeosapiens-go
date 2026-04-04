-- +goose Up
-- +goose StatementBegin
delete from orders;
alter table orders
add column order_number bigint not null,
add column billing_tax_id bytea;
create unique index on orders (order_number);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table orders drop column order_number, drop column billing_tax_id;
-- +goose StatementEnd
