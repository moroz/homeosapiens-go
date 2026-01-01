-- name: UpsertAsset :one
insert into assets (id, object_key, original_filename) values ($1, $2, $3)
on conflict (id) do update set object_key = excluded.object_key, original_filename = excluded.original_filename, updated_at = now()
returning *;

-- name: UpsertHost :one
INSERT INTO hosts (id, salutation, given_name, family_name, profile_picture_id, country)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (id) do nothing
returning *;

-- name: UpsertVenue :one
INSERT INTO venues (id, name_en, name_pl, street, city_en, city_pl, postal_code, country_code)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (id) DO UPDATE SET
    name_en = excluded.name_en,
    name_pl = excluded.name_pl,
    street = excluded.street,
    city_en = excluded.city_en,
    city_pl = excluded.city_pl,
    postal_code = excluded.postal_code,
    country_code = excluded.country_code,
    updated_at = now()
returning *;

-- name: UpsertEvent :one
INSERT INTO events (id, event_type, title_en, title_pl, slug, starts_at, ends_at, is_virtual, description_en, description_pl, venue_id, base_price_amount, base_price_currency)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
ON CONFLICT (id) DO UPDATE SET
    event_type = excluded.event_type,
    title_en = excluded.title_en,
    title_pl = excluded.title_pl,
    slug = excluded.slug,
    starts_at = excluded.starts_at,
    ends_at = excluded.ends_at,
    is_virtual = excluded.is_virtual,
    description_en = excluded.description_en,
    description_pl = excluded.description_pl,
    venue_id = excluded.venue_id,
    base_price_amount = excluded.base_price_amount,
    base_price_currency = excluded.base_price_currency,
    updated_at = now()
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

