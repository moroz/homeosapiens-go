-- name: InsertProduct :one
insert into products (product_type, title_pl, title_en, base_price_amount, base_price_currency) values ($1, $2, $3, $4, $5) returning *;