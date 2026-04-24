-- name: ListAssetsByIDs :many
select * from assets where id = any(@asset_ids::uuid[]) order by id;