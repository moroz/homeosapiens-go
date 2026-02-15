-- name: GetEventById :one
select * from events where id = (@id::text)::uuid;

-- name: GetEventBySlug :one
select * from events where slug = @slug::text;

-- name: ListEvents :many
select e.id, e.slug, e.title_en, e.title_pl, e.is_virtual, e.base_price_amount, e.base_price_currency,
       e.event_type, e.starts_at, e.ends_at, e.subtitle_pl, e.subtitle_en,
       e.venue_street, e.venue_city_en, e.venue_city_pl, e.venue_country_code
from events e
order by e.starts_at desc;

-- name: ListHostsForEvents :many
select eh.event_id, h.*, a.object_key profile_picture_url
from hosts h
join events_hosts eh on eh.host_id = h.id
left join assets a on h.profile_picture_id = a.id
where eh.event_id = any(@EventIDs::uuid[])
order by eh.host_id, eh.position;

-- name: ListPricesForEvents :many
select p.* from event_prices p
where p.event_id = any(@EventIDs::uuid[])
order by p.event_id, p.priority;

-- name: ListEventRegistrationsForUserForEvents :many
select er.* from event_registrations er
where er.event_id = any(@EventIDs::uuid[])
and er.user_id = @UserID::uuid;

-- name: GetFreeEventById :one
select * from events where (base_price_amount is null or base_price_amount = 0) and id = (@id::text)::uuid;

-- name: GetPaidEventById :one
select * from events where base_price_amount is not null and base_price_amount > 0 and id = (@id::text)::uuid;