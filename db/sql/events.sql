-- name: GetEventById :one
select * from events where id = $1;

-- name: GetEventBySlug :one
select * from events where slug = $1;

-- name: ListEvents :many
select e.id, e.slug, e.title_en, e.title_pl, e.is_virtual, p.base_price_amount, p.base_price_currency,
       e.event_type, e.starts_at, e.ends_at, e.subtitle_pl, e.subtitle_en,
       e.venue_street, e.venue_city_en, e.venue_city_pl, e.venue_country_code
from events e
left join products p on e.product_id = p.id
order by e.starts_at desc;

-- name: PaginateEvents :many
select * from events
order by starts_at desc
limit (@per_page::int) offset (((@page::int) - 1) * @per_page::int);

-- name: CountEvents :one
select count(*) from events;

-- name: ListProductsForEvents :many
select e.id event_id, sqlc.embed(p)
from events e
join products p on e.product_id = p.id
where e.id = any(@EventIDs::uuid[])
order by 1;

-- name: ListHostsForEvents :many
select eh.event_id, h.*, a.object_key profile_picture_url
from hosts h
join events_hosts eh on eh.host_id = h.id
left join assets a on h.profile_picture_id = a.id
where eh.event_id = any(@EventIDs::uuid[])
order by eh.host_id, eh.position;

-- name: ListPricesForEvents :many
select e.id event_id, sqlc.embed(p) from product_prices p
join events e on e.product_id = p.id
where e.id = any(@EventIDs::uuid[])
order by e.id, p.priority;

-- name: ListEventRegistrationsForUserForEvents :many
select er.* from event_registrations er
where er.event_id = any(@EventIDs::uuid[])
and er.user_id = @UserID::uuid;

-- name: GetFreeEventById :one
select * from events where product_id is null and id = $1;

-- name: GetPaidEventById :one
select sqlc.embed(e), sqlc.embed(p)
from events e
join products p on e.product_id = p.id
where e.id = $1;

-- name: UpdateEvent :one
update events
set title_pl = $1, title_en = $2, subtitle_pl = $3, subtitle_en = $4, description_pl = $5, description_en = $6
where id = @event_id::uuid
returning *;

-- name: InsertEvent :one
insert into events (title_en, title_pl, starts_at, ends_at, is_virtual, description_en, description_pl, event_type, slug, subtitle_en, subtitle_pl, venue_name_en, venue_name_pl, venue_street, venue_city_en, venue_city_pl, venue_postal_code, venue_country_code, product_id)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12,$13,$14, $15, $16, $17, $18, $19)
returning *;

-- name: InsertEventHost :one
insert into events_hosts (event_id, host_id, position) values ($1, $2, $3) returning *;