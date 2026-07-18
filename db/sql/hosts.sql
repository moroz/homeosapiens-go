-- name: ListHosts :many
select * from hosts
order by family_name, given_name;
