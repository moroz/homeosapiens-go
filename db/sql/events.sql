-- name: GetEventById :one
select * from events where id = $1;

-- name: ListEvents :many
select e.id, e.slug, e.title_en, e.title_pl, e.is_virtual, e.base_price_amount, e.base_price_currency,
       e.venue_id, e.event_type, e.starts_at, e.ends_at,
       v.street venue_street, v.city_en venue_city_en, v.city_pl venue_city_pl, v.country_code venue_country_code
from events e
left join venues v on e.venue_id = v.id
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

-- name: GetEventBySlug :one
select e.*
from events e
left join venues v on e.venue_id = v.id
where e.slug = @slug::text;
