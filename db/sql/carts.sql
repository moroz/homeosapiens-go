-- name: InsertCartLineItem :one
insert into cart_line_items as cli (cart_id, product_id, quantity) values ($1, $2, 1) on conflict (cart_id, product_id)
do update set quantity = cli.quantity + excluded.quantity returning *;

-- name: GetCart :one
select cli.cart_id, count(cli.id) item_count, sum(cli.quantity * p.base_price_amount)::decimal product_total
from cart_line_items cli
join products p on cli.product_id = p.id
where cli.cart_id = @cart_id::uuid
group by 1;

-- name: CountCartLineItemQuantitiesForProducts :many
select e.id event_id, c.quantity
from cart_line_items c
join events e on c.product_id = e.product_id
where e.id = any(@event_ids::uuid[]) and c.cart_id = @cart_id::uuid;

-- name: GetCartItemsByCartId :many
select c.*, (p.base_price_amount * c.quantity)::decimal as subtotal, p.base_price_amount, p.title_en, p.title_pl, e.slug::text slug
from cart_line_items c
join products p on c.product_id = p.id
left join events e on e.product_id = p.id
where c.cart_id = @cart_id::uuid and (e.id is not null);

-- name: DeleteCartItem :one
delete from cart_line_items cli where cart_id = @cart_id::uuid and product_id = @product_id::uuid returning id;

-- name: DeleteCart :exec
delete from cart_line_items where cart_id = @cart_id::uuid;