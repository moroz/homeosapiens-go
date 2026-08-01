-- name: InsertProduct :one
insert into products (product_type, title_pl, title_en, base_price_amount, base_price_currency) values ($1, $2, $3, $4, $5) returning *;

-- name: GetProductById :one
select * from products where id = $1;

-- name: UpdateProductPrice :one
update products set base_price_amount = $1, base_price_currency = $2, updated_at = now()
where id = @product_id::uuid
returning *;
