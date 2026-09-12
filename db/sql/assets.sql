-- name: ListAssetsByIDs :many
select * from assets where id = any(@asset_ids::uuid[]) order by id;

-- name: InsertAsset :one
insert into assets (object_key, original_filename) values ($1, $2) returning *;
