-- name: InsertEventRegistration :one
insert into event_registrations (event_id, user_id) values ($1, $2)
on conflict (event_id, user_id) do nothing
returning *;

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

-- name: CountRegistrationsForEvents :many
select er.event_id, count(er.id) from event_registrations er
where er.event_id = any(@EventIDs::uuid[])
group by 1;