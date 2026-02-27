-- name: ListOrders :many
select * from orders order by id;

-- name: InsertOrder :one
insert into orders (user_id, grand_total, currency, billing_given_name_encrypted, billing_family_name_encrypted, billing_phone_encrypted, billing_street_encrypted, billing_house_number_encrypted, billing_apartment_number_encrypted, billing_city_encrypted, billing_postal_code_encrypted, billing_country, email_encrypted)
values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) returning *;