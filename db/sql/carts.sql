-- name: InsertCart :one
insert into carts (owner_id) values ($1)
on conflict (owner_id) where owner_id is not null do update set updated_at = now() at time zone 'utc'
returning *;

-- name: InsertCartLineItem :one
insert into cart_line_items as cli (cart_id, event_id, quantity) values ($1, $2, 1) on conflict (cart_id, event_id)
do update set quantity = cli.quantity + excluded.quantity returning *;

-- name: GetCart :one
select sqlc.embed(c), count(cli.id) item_count, sum(cli.quantity * e.base_price_amount)::decimal product_total
from carts c
join cart_line_items cli on cli.cart_id = c.id
join events e on cli.event_id = e.id
where
    (sqlc.narg('cart_id')::uuid is not null and c.id = sqlc.narg('cart_id')::uuid)
    or
    (sqlc.narg('cart_id')::uuid is null and sqlc.narg('owner_id')::uuid is not null and c.owner_id = sqlc.narg('owner_id')::uuid)
    group by c.id, c.owner_id, c.inserted_at, c.updated_at;

-- name: CountCartLineItemQuantitiesForEvents :many
select c.event_id, c.quantity
from cart_line_items c
where c.event_id = any(@event_ids::uuid[]) and c.cart_id = @cart_id::uuid;