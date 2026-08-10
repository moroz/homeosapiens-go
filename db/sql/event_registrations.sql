-- name: InsertEventRegistration :one
-- The ON CONFLICT DO UPDATE clause is required to force Postgres to return the existing row.
-- Since the event_id column is the same for all insert of the same conflicting event, 
insert into event_registrations as er (event_id, user_id) values ($1, $2)
on conflict (event_id, user_id) do update set event_id = excluded.event_id
returning sqlc.embed(er), (xmax = 0)::boolean new_record;

-- name: DeleteEventRegistration :one
delete from event_registrations where event_id = $1 and user_id = $2 returning id;

-- name: GetLastEventRegistration :one
-- Only for testing
select * from event_registrations order by id desc limit 1;

-- name: GetLastEventRegistrationWithDetails :one
-- For dev email preview only
select sqlc.embed(u), sqlc.embed(e)
from event_registrations er
join users u on u.id = er.user_id
join events e on e.id = er.event_id
order by er.inserted_at desc
limit 1;

-- name: ClaimEventRegistrationsForReminder24h :many
-- Stamps and returns the registrations whose day-ahead reminder is due, so that
-- a second worker running concurrently cannot pick the same one up. The stamps
-- live on the registration rather than on the event, so that somebody
-- registering after the scan has run still gets whatever reminders are left.
-- Events starting within the hour are left to the one-hour reminder, which keeps
-- a late registration from firing both reminders at once.
update event_registrations set reminder_24h_sent_at = now()
where id in (
  select er.id from event_registrations er
  join events e on e.id = er.event_id
  where e.published_at is not null
    and er.reminder_24h_sent_at is null
    and e.starts_at > now() + interval '1 hour'
    and e.starts_at <= now() + interval '24 hours'
  for update of er skip locked
)
returning event_id, user_id;

-- name: ClaimEventRegistrationsForReminder1h :many
-- The one-hour counterpart of ClaimEventRegistrationsForReminder24h.
update event_registrations set reminder_1h_sent_at = now()
where id in (
  select er.id from event_registrations er
  join events e on e.id = er.event_id
  where e.published_at is not null
    and er.reminder_1h_sent_at is null
    and e.starts_at > now()
    and e.starts_at <= now() + interval '1 hour'
  for update of er skip locked
)
returning event_id, user_id;

-- name: CountRegistrationsForEvents :many
select er.event_id, count(er.id) from event_registrations er
where er.event_id = any(@EventIDs::uuid[])
group by 1;

-- name: PaginateEventRegistrations :many
select distinct on (er.id)
  u.id, o.id order_id, o.order_number, u.given_name_encrypted, u.family_name_encrypted,
  u.email_encrypted, er.inserted_at
from event_registrations er
join users u on er.user_id = u.id
join events e on er.event_id = e.id
left join order_line_items oli on oli.product_id = e.product_id
left join orders o on o.id = oli.order_id and o.user_id = er.user_id and o.paid_at is not null and o.cancelled_at is null
where er.event_id = @event_id::uuid
order by er.id desc, o.paid_at desc
limit (@per_page::int) offset (((@page::int) - 1) * @per_page::int);

-- name: CountEventRegistrations :one
select count(er.id) from event_registrations er
where er.event_id = @event_id::uuid;

-- name: ListEligibleUsersForEvent :many
select u.id, u.given_name_encrypted, u.family_name_encrypted, u.email_encrypted,
  (er.id is null)::boolean as can_register
from users u
left join event_registrations er on er.user_id = u.id and er.event_id = @event_id::uuid
order by er.id desc;
