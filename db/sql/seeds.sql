-- name: UpsertAsset :one
insert into assets (id, object_key, original_filename) values ($1, $2, $3)
on conflict (id) do update set object_key = excluded.object_key, original_filename = excluded.original_filename, updated_at = now()
returning *;

-- name: UpsertHost :one
INSERT INTO hosts (id, salutation, given_name, family_name, profile_picture_id, country)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (id) do nothing
returning *;

-- name: UpsertEvent :one
INSERT INTO events (id, event_type, title_en, title_pl, slug, starts_at, ends_at, is_virtual, description_en, description_pl, base_price_amount, base_price_currency, subtitle_en, subtitle_pl, venue_street, venue_city_en, venue_city_pl, venue_name_en, venue_name_pl, venue_country_code, venue_postal_code)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)
returning *;

-- name: UpsertEventHost :one
INSERT INTO events_hosts (event_id, host_id, position)
VALUES ($1, $2, $3)
ON CONFLICT (event_id, host_id) DO UPDATE SET
    position = excluded.position,
    updated_at = now()
returning *;

-- name: UpsertEventPrice :one
INSERT INTO event_prices (event_id, price_type, rule_type, price_amount, price_currency, discount_code, priority, is_active, valid_from, valid_until)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
on conflict do nothing
returning *;

-- name: UpsertVideo :one
INSERT INTO videos (id, event_id, provider, title_en, title_pl, slug, is_public)
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (id) DO UPDATE SET
    event_id = excluded.event_id,
    provider = excluded.provider,
    title_en = excluded.title_en,
    title_pl = excluded.title_pl,
    slug = excluded.slug,
    is_public = excluded.is_public,
    updated_at = now()
returning *;

-- name: UpsertVideoSource :one
INSERT INTO video_sources (id, video_id, content_type, codec, object_key)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (id) DO UPDATE SET
    video_id = excluded.video_id,
    content_type = excluded.content_type,
    codec = excluded.codec,
    object_key = excluded.object_key,
    updated_at = now()
returning *;

