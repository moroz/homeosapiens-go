-- name: GetUserByEmail :one
select * from users where email_hash = $1;

-- name: GetUserByAccessToken :one
select u.* from user_tokens ut
join users u on ut.user_id = u.id
where ut.valid_until > (now() at time zone 'utc')
and ut.token = @token::bytea and ut.context = 'access';

-- name: InsertUser :one
insert into users (email_encrypted, email_hash, salutation, given_name_encrypted, family_name_encrypted, country, profession, organization, company, password_hash) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) returning *;

-- name: InsertUserToken :one
insert into user_tokens (user_id, context, token, valid_until) values ($1, $2, $3, $4) returning *;

-- name: DeleteUserToken :one
delete from user_tokens where token = $1 returning true;

-- name: FindOrCreateUserFromClaims :one
insert into users (email_encrypted, email_hash, given_name_encrypted, family_name_encrypted, profile_picture)
values ($1, $2, $3, $4, $5)
on conflict (email_hash) do update
set given_name_encrypted = excluded.given_name_encrypted, family_name_encrypted = excluded.family_name_encrypted, profile_picture = excluded.profile_picture, updated_at = now() at time zone 'utc'
returning *;
