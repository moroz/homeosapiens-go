-- name: InsertEventRegistration :one
insert into event_registrations (event_id, user_id) values ($1, $2) returning *;

-- name: DeleteEventRegistration :one
delete from event_registrations where event_id = $1 and user_id = $2 returning id;

-- name: GetLastEventRegistration :one
-- Only for testing
select * from event_registrations order by id desc limit 1;

-- name: CountRegistrationsForEvents :many
select er.event_id, count(er.id) from event_registrations er
where er.event_id = any(@EventIDs::uuid[])
group by 1;