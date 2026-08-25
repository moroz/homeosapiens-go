-- name: ListHosts :many
select * from hosts
order by family_name, given_name;

-- name: PaginateHosts :many
select * from hosts
order by family_name, given_name
limit (@per_page::int) offset (((@page::int) - 1) * @per_page::int);

-- name: CountHosts :one
select count(*) from hosts;

-- name: GetHostById :one
select * from hosts where id = $1;

-- name: InsertHost :one
insert into hosts (salutation, given_name, family_name, country)
values ($1, $2, $3, $4)
returning *;

-- name: UpdateHost :one
update hosts set salutation = $2, given_name = $3, family_name = $4, country = $5, updated_at = now()
where id = $1
returning *;

-- name: DeleteHost :exec
delete from hosts where id = $1;
