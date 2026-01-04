-- name: InsertEventRegistration :one
insert into event_registrations (event_id, user_id, is_host, given_name_encrypted, family_name_encrypted, email_encrypted, country, attending_in_person, email_confirmed_at)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
returning *;
