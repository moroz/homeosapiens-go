-- name: ListHosts :many
select * from hosts
order by family_name, given_name;

-- name: PaginateHosts :many
select * from hosts
order by family_name, given_name
limit @per_page offset ((@page - 1) * @per_page);

-- name: CountHosts :one
select count(*) from hosts;
