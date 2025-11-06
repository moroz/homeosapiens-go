-- name: ListEvents :many
select * from events order by starts_at desc;