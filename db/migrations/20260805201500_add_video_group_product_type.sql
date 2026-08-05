-- +goose Up
-- Video groups are sold on their own, so they need a product type of their own.
alter type product_type add value if not exists 'video_group';

-- +goose Down
-- Enum values cannot be dropped in PostgreSQL; the type would have to be
-- recreated, which is not worth it for an additive change.
select 1;
