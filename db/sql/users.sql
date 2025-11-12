-- name: GetUserByEmail :one
select * from users where email = $1;

-- name: GetUserByAccessToken :one
select u.* from user_tokens ut
join users u on ut.user_id = u.id
where ut.valid_until > (now() at time zone 'utc')
and ut.token = @token::bytea and ut.context = 'access';

-- name: InsertUser :one
insert into users (email, salutation, given_name, family_name, country, profession, organization, company, password_hash) values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning *;

-- name: InsertUserToken :one
insert into user_tokens (user_id, context, token, valid_until) values ($1, $2, $3, $4) returning *;