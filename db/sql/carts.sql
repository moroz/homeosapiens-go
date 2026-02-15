-- name: InsertCart :one
insert into carts (owner_id) values ($1) returning *;

-- name: InsertCartLineItem :one
insert into cart_line_items (cart_id, event_id) values ($1, $2) on conflict (cart_id, event_id) do update set quantity = quantity + excluded.quantity returning *;

-- name: GetCart :one
select c.id, count(cli.id), sum(cli.quantity * e.base_price_amount)
from carts c
join cart_line_items cli on cli.cart_id = c.id
join events e on cli.event_id = e.id
where c.id = (@cart_id::text)::uuid
group by 1;