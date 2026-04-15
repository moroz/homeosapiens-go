-- name: InsertCartLineItem :one
insert into cart_line_items as cli (cart_id, product_id, quantity) values ($1, $2, 1) on conflict (cart_id, product_id)
do update set quantity = cli.quantity + excluded.quantity returning *;

-- name: GetCart :one
select cli.cart_id, count(cli.id) item_count, sum(cli.quantity * e.base_price_amount)::decimal product_total
from cart_line_items cli
join events e on cli.product_id = e.id
where cli.cart_id = @cart_id::uuid 
group by 1;

-- name: CountCartLineItemQuantitiesForProducts :many
select c.product_id, c.quantity
from cart_line_items c
where c.product_id = any(@product_ids::uuid[]) and c.cart_id = @cart_id::uuid;

-- name: GetCartItemsByCartId :many
select sqlc.embed(c), sqlc.embed(p), (p.base_price_amount * c.quantity)::decimal as subtotal
from cart_line_items c
join products p on c.product_id = p.id
where c.cart_id = @cart_id::uuid;

-- name: DeleteCartItem :one
delete from cart_line_items cli where cart_id = @cart_id::uuid and product_id = @product_id::uuid returning id;

-- name: DeleteCart :exec
delete from cart_line_items where cart_id = @cart_id::uuid;