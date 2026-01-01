-- name: UpsertAsset :one
insert into assets (id, object_key, original_filename) values ($1, $2, $3)
on conflict (id) do update set object_key = excluded.object_key, original_filename = excluded.original_filename, updated_at = now()
returning *;

-- name: UpsertHost :one
INSERT INTO hosts (id, salutation, given_name, family_name, profile_picture_id, country)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (id) do nothing
returning *;
