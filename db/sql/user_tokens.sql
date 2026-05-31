-- name: InsertUserToken :one
insert into user_tokens (user_id, context, token, valid_until) values ($1, $2, $3, $4) returning *;

-- name: DeleteUserToken :one
delete from user_tokens where token = $1 returning true;

-- name: DeletePreexistingEmailVerificationTokens :exec
delete from user_tokens where user_id = $1 and context = 'email_verification';

-- name: FindUserByUserToken :one
select u.* from users u
join user_tokens ut on u.id = ut.user_id
where ut.valid_until > now()
and ut.token = @token and ut.context = @context;

-- name: FindUserAndTokenByUserToken :one
select sqlc.embed(u), sqlc.embed(ut) from users u
join user_tokens ut on u.id = ut.user_id
where ut.valid_until > now()
and ut.token = @token and ut.context = @context;

-- name: VacuumUserTokens :exec
delete from user_tokens where valid_until < now();