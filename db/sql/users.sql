-- name: GetUserByEmail :one
select * from users where email_hash = $1;

-- name: GetUserByAccessToken :one
select u.* from user_tokens ut
join users u on ut.user_id = u.id
where ut.valid_until > now()
and ut.token = @token and ut.context = @context;

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

-- name: PaginateUsers :many
select * from users order by id
limit (@per_page::int) offset (((@page::int) - 1) * @per_page::int);

-- name: CountUsers :one
select count(*) from users;

-- name: SetUserLastLogin :exec
update users set last_login_ip = $1, last_login_at = now(), updated_at = now()
where id = $2;

-- name: UpdateUserPreferredLocale :exec
update users set preferred_locale = $1, updated_at = now() where id = $2;

-- name: UpdateUserPreferredTimezone :exec
update users set preferred_timezone_encrypted = $1, updated_at = now() where id = $2;

-- name: GetUserByID :one
select * from users where id = $1;

-- name: VerifyEmailAddressByUserToken :one
update users u set email_confirmed_at = now(), updated_at = now()
from user_tokens ut
where ut.token = $1 and ut.valid_until > now() and ut.user_id = u.id and u.email_confirmed_at is null
returning u.*;

-- name: UpdateUserPassword :one
update users
set password_hash = $1, email_confirmed_at = coalesce(email_confirmed_at, now()), updated_at = now()
where id = $2
returning *;

-- name: ListUserProductAccess :many
select p.id, p.product_type, p.title_pl, p.title_en, upa.order_id, upa.inserted_at granted_at,
gb.id granted_by_user_id, gb.given_name_encrypted granted_by_given_name, gb.family_name_encrypted granted_by_family_name
from user_product_access upa
join products p on p.id = upa.product_id
left join users gb on gb.id = upa.granted_by_user_id
where upa.user_id = $1
order by upa.id desc;

-- name: GrantProductAccess :exec
-- Grants access outside of an order, e.g. by hand from the admin panel or from
-- an import script. granted_by_user_id records the administrator who did it.
insert into user_product_access (user_id, product_id, granted_by_user_id, inserted_at)
values ($1, $2, sqlc.narg(granted_by_user_id), coalesce(sqlc.narg(inserted_at)::timestamp, now()))
on conflict (user_id, product_id) do nothing;

-- name: ListEventRegistrationsByUserID :many
select e.id, e.slug, e.title_pl, e.title_en, e.starts_at, e.ends_at, er.inserted_at registered_at
from event_registrations er
join events e on e.id = er.event_id
where er.user_id = $1
order by e.starts_at desc;
