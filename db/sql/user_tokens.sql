-- name: InsertUserToken :one
insert into user_tokens (user_id, context, token, valid_until) values ($1, $2, $3, $4) returning *;

-- name: DeleteUserToken :one
delete from user_tokens where token = $1 returning true;
