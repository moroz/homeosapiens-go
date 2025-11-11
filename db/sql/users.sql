-- name: GetUserByEmail :one
select * from users where email = $1;

-- name: GetUserByAccessToken :one
select u.* from user_tokens ut
join users u on ut.user_id = u.id
where ut.valid_until > (now() at time zone 'utc')
and ut.token = @token::bytea and ut.context = 'access';