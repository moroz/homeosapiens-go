-- name: GetEventById :one
select * from events where id = $1;

-- name: GetEventBySlug :one
select * from events where slug = $1;

-- name: ListPublishedEvents :many
select e.id, e.slug, e.title_en, e.title_pl, e.is_virtual, p.base_price_amount, p.base_price_currency,
       e.event_type, e.starts_at, e.ends_at, e.subtitle_pl, e.subtitle_en,
       e.venue_street, e.venue_city_en, e.venue_city_pl, e.venue_country_code
from events e
left join products p on e.product_id = p.id
where published_at is not null
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
order by eh.event_id, eh.position;

-- name: ListPricesForEvents :many
select e.id event_id, sqlc.embed(p) from product_prices p
join events e on e.product_id = p.id
where e.id = any(@EventIDs::uuid[])
order by e.id, p.priority;

-- name: ListEventRegistrationsForUserForEvents :many
select er.* from event_registrations er
where er.event_id = any(@EventIDs::uuid[])
and er.user_id = @UserID::uuid;

-- name: GetRegisterableFreeEventById :one
-- Registration is only possible while the event is still running: once it has
-- ended, there is nothing left to sign up for.
select * from events where product_id is null and id = $1 and ends_at > now();

-- name: GetPaidEventById :one
select sqlc.embed(e), sqlc.embed(p)
from events e
join products p on e.product_id = p.id
where e.id = $1;

-- name: InsertEvent :one
insert into events (title_en, title_pl, starts_at, ends_at, is_virtual, description_en, description_pl, event_type, slug, subtitle_en, subtitle_pl, venue_name_en, venue_name_pl, venue_street, venue_city_en, venue_city_pl, venue_postal_code, venue_country_code, product_id, meeting_url)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12,$13,$14, $15, $16, $17, $18, $19, $20)
returning *;

-- name: InsertEventHost :one
insert into events_hosts (event_id, host_id, position) values ($1, $2, $3) returning *;

-- name: DeleteEventHosts :exec
delete from events_hosts where event_id = $1;

-- name: PublishEvent :one
update events set published_at = now(), updated_at = now()
where id = $1 and published_at is null returning *;

-- name: UnpublishEvent :one
update events set published_at = null, updated_at = now()
where id = $1 returning *;

-- name: DeleteEvent :execrows
delete from events where id = $1;

-- name: CountRegistrationsForEvent :one
select count(*) from event_registrations where event_id = $1;

-- name: ClaimEventsForReminder24h :many
-- Stamps and returns the events whose day-ahead reminder is due, so that a
-- second worker running concurrently cannot pick the same event up. Events
-- starting within the hour are left to the one-hour reminder, which keeps an
-- event published at the last minute from firing both reminders at once.
update events set reminder_24h_sent_at = now()
where id in (
  select id from events
  where published_at is not null
    and reminder_24h_sent_at is null
    and starts_at > now() + interval '1 hour'
    and starts_at <= now() + interval '24 hours'
  for update skip locked
)
returning id;

-- name: ClaimEventsForReminder1h :many
-- The one-hour counterpart of ClaimEventsForReminder24h.
update events set reminder_1h_sent_at = now()
where id in (
  select id from events
  where published_at is not null
    and reminder_1h_sent_at is null
    and starts_at > now()
    and starts_at <= now() + interval '1 hour'
  for update skip locked
)
returning id;
