-- name: ListOrders :many
select * from orders order by id;

-- name: InsertOrder :one
insert into orders (user_id, grand_total, currency, billing_given_name_encrypted, billing_family_name_encrypted, billing_phone_encrypted, billing_city_encrypted, billing_postal_code_encrypted, billing_country, email_encrypted, billing_address_line1_encrypted, billing_address_line2_encrypted)
values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) returning *;

-- name: InsertOrderLineItem :one
insert into order_line_items (order_id, event_id, event_title, event_price_amount) VALUES ($1, $2, $3, $4) returning *;

-- name: StoreCheckoutSessionIDOnOrder :one
update orders set stripe_checkout_session_id = $1, updated_at = now() where id = $2 returning *;