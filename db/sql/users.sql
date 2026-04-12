-- name: GetUserByEmail :one
select * from users where email_hash = $1;

-- name: GetUserByAccessToken :one
select u.* from user_tokens ut
join users u on ut.user_id = u.id
where ut.valid_until > now()
and ut.token = @token and ut.context = 'access';

-- name: InsertUser :one
insert into users (email_encrypted, email_hash, salutation, given_name_encrypted, family_name_encrypted, country, profession, organization, company, password_hash, preferred_locale) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) returning *;

-- name: UpsertUserFromSeedData :one
insert into users (email_encrypted, email_hash, given_name_encrypted, family_name_encrypted, country, password_hash, user_role, email_confirmed_at, preferred_locale)
values ($1, $2, $3, $4, $5, $6, coalesce(sqlc.narg(user_role)::text::user_role, 'Regular'), $7, coalesce(sqlc.narg(preferred_locale)::text::locale, 'pl'))
on conflict (email_hash) do update set updated_at = now()
returning *;

-- name: FindOrCreateUserFromClaims :one
insert into users (email_encrypted, email_hash, given_name_encrypted, family_name_encrypted, profile_picture, preferred_locale, email_confirmed_at, google_oauth_last_used_at)
values ($1, $2, $3, $4, $5, $6, now(), now())
on conflict (email_hash) do update
set given_name_encrypted = excluded.given_name_encrypted, family_name_encrypted = excluded.family_name_encrypted, profile_picture = excluded.profile_picture, updated_at = now(), email_confirmed_at = coalesce(users.email_confirmed_at, excluded.email_confirmed_at), preferred_locale = coalesce(users.preferred_locale, excluded.preferred_locale), google_oauth_last_used_at = now()
returning *;

-- name: UpdateUserProfile :one
update users
set given_name_encrypted = $1, family_name_encrypted = $2, profession = $3, licence_number_encrypted = $4, country = $5, updated_at = now()
where id = $6 returning *;

-- name: ListUsers :many
select * from users order by id;

-- name: SetUserLastLogin :exec
update users set last_login_ip = $1, last_login_at = now(), updated_at = now()
where id = $2;

-- name: UpdateUserPreferredLocale :exec
update users set preferred_locale = $1 where id = $2;

-- name: GetUserByID :one
select * from users where id = $1;

-- name: VerifyEmailAddressByUserToken :one
update users u set email_confirmed_at = now(), updated_at = now()
from user_tokens ut
where ut.token = $1 and ut.valid_until > now() and ut.user_id = u.id and u.email_confirmed_at is null
returning u.*;
