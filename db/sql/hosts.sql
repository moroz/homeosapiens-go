-- name: ListHosts :many
select * from hosts
order by family_name, given_name;

-- name: PaginateHosts :many
select * from hosts
order by family_name, given_name
limit (@per_page::int) offset (((@page::int) - 1) * @per_page::int);

-- name: CountHosts :one
select count(*) from hosts;
