-- name: ListOrders :many
select * from orders order by id;

-- name: InsertOrder :one
insert into orders (order_number, user_id, grand_total, currency, billing_given_name_encrypted, billing_family_name_encrypted, billing_phone_encrypted, billing_city_encrypted, billing_postal_code_encrypted, billing_country, email_encrypted, billing_address_line1_encrypted, billing_address_line2_encrypted, billing_tax_id)
values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14) returning *;

-- name: InsertOrderLineItem :one
insert into order_line_items (order_id, event_id, event_title, event_price_amount) VALUES ($1, $2, $3, $4) returning *;

-- name: StoreCheckoutSessionIDOnOrder :one
update orders set stripe_checkout_session_id = $1, updated_at = now() where id = $2 returning *;

-- name: GetLastOrderID :one
select id from orders order by id desc limit 1;

-- name: GetOrderByID :one
select * from orders where id = $1;

-- name: GetOrderLineItemsForOrderID :many
select * from order_line_items where order_id = $1 order by id;

-- name: GetOrderByCheckoutSessionID :one
select * from orders where stripe_checkout_session_id = @session_id::text;

-- name: GetOrderByCheckoutSessionIDForUpdate :one
select * from orders where stripe_checkout_session_id = @session_id::text for update;

-- name: MarkOrderAsPaid :one
update orders set paid_at = now(), updated_at = now() where id = $1 returning *;

-- name: GenerateNextOrderNumberForDate :one
WITH prefix AS (
    SELECT (
        (EXTRACT(YEAR FROM sqlc.arg(date)::date) % 100) * 100_000_000 -- YY
        + EXTRACT(MONTH FROM sqlc.arg(date)::date)      * 1_000_000   -- MM
        + EXTRACT(DAY FROM sqlc.arg(date)::date)        * 10_000      -- DD
    )::bigint AS val
)
SELECT ((COALESCE(MAX(order_number) % 10000, 0) + 1) + prefix.val)::bigint
FROM prefix
LEFT JOIN orders ON order_number BETWEEN prefix.val and prefix.val + 9999
GROUP BY prefix.val;
