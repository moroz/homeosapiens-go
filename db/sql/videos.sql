-- name: ListVideos :many
select v.* from videos v
order by v.id desc;

-- name: ListVideosForUser :many
select v.* from videos v
where v.is_public = true
   or v.id in (select v.id from videos v join event_registrations er on v.event_id = er.event_id
               where er.user_id = @UserID::uuid)
   or exists (select 1 from users u where u.id = @UserID::uuid and u.user_role = 'Administrator');

select v.* from videos v
left join events e on v.event_id = e.id
left join event_registrations er on er.event_id = e.id
where er.user_id = $1 or v.is_public = true
order by v.id desc;

-- name: ListPublicVideos :many
select * from videos v
where is_public = true
order by v.id desc;

-- name: ListVideoSourcesForVideos :many
select * from video_sources vs
where vs.video_id = any(@VideoIDs::uuid[])
order by vs.video_id, vs.id;
