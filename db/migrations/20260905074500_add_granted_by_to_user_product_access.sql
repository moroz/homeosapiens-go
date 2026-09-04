-- +goose Up
alter table user_product_access
add column granted_by_user_id uuid references users (id) on delete set null;

comment on column user_product_access.granted_by_user_id is
  'The administrator who granted the access by hand. Null for purchases and imports.';

-- +goose Down
alter table user_product_access drop column granted_by_user_id;
