-- name: ListEvents :many
select e.*, v.*
from events e
left join venues v on e.venue_id = v.id
order by e.starts_at desc;

-- name: ListHostsForEvents :many
select h.*, a.object_key profile_picture_url
from hosts h
join events_hosts eh on eh.host_id = h.id
left join assets a on h.profile_picture_id = a.id
where eh.event_id = any(@EventIDs)
order by eh.host_id, eh.position;