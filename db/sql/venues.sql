-- name: GetVenueById :one
select * from venues where id = $1;