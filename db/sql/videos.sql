-- name: ListVideos :many
select v.* from videos v
order by v.id desc;

-- name: ListVideoGroupsForUser :many
select sqlc.embed(vg), a.has_access from video_groups vg
join user_video_group_access a on vg.id = a.video_group_id and a.user_id = @user_id::uuid
order by vg.id desc;

-- name: GetMinMaxRecordedDatesForVideoGroups :many
select vg.id, min(v.recorded_on)::date min_recorded_on, max(v.recorded_on):: date max_recorded_on
from video_groups vg
join video_groups_videos vgv on vg.id = vgv.video_group_id
join videos v on v.id = vgv.video_id
where vg.id = any(@video_group_ids::uuid[])
group by 1
order by 1;

-- name: GetVideoGroupForUserBySlug :one
select sqlc.embed(vg), a.has_access from video_groups vg
join user_video_group_access a on vg.id = a.video_group_id and a.user_id = @user_id::uuid
where (sqlc.narg(slug)::text is null or vg.slug = sqlc.narg(slug))
limit 1;

-- name: GetVideoForUser :one
select sqlc.embed(v), a.has_access from videos v
join video_groups_videos vgv on vgv.video_id = v.id
join video_groups vg on vgv.video_group_id = vg.id
join user_video_group_access a on vg.id = a.video_group_id and a.user_id = @user_id::uuid
where v.slug = @video_slug and vg.slug = @group_slug
limit 1;

-- name: ListVideosForVideoGroup :many
select v.* from videos v
join video_groups_videos vgv on vgv.video_id = v.id
where vgv.video_group_id = $1
order by position;

-- name: InsertVideoGroup :one
insert into video_groups (title_en, title_pl, slug, product_id) values ($1, $2, $3, $4) returning *;

-- name: UpsertVideoGroup :one
insert into video_groups (id, title_en, title_pl, slug, product_id) values ($1, $2, $3, $4, $5)
on conflict (slug) do update set title_en = excluded.title_en, title_pl = excluded.title_pl, product_id = excluded.product_id
returning *;

-- name: ListVideoSourcesForVideos :many
select * from video_sources vs
where vs.video_id = any(@video_ids::uuid[])
order by vs.video_id, vs.priority;

-- name: AddVideoToVideoGroup :one
insert into video_groups_videos (video_id, video_group_id, position)
select $1, $2, coalesce(max(position), -1) + 1
from video_groups_videos where video_group_id = $2
on conflict (video_id, video_group_id) do nothing
returning *;

-- name: GetVideoThumbnailData :one
select sqlc.embed(v), h.given_name host_given_name, h.family_name host_family_name, a.object_key host_profile_picture_url
from videos v
left join hosts h on h.id = v.host_id
left join assets a on a.id = h.profile_picture_id
where v.id = @id::uuid;

-- name: InsertVideo :one
insert into videos (provider, title_en, title_pl, slug, duration_seconds, recorded_on, host_id, thumbnail_en_id, thumbnail_pl_id)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
returning *;

-- name: ListYoutubeVideos :many
select * from videos where provider = 'youtube' order by recorded_on desc;

-- name: ListHostsForVideos :many
select vh.video_id, sqlc.embed(h)
from video_hosts vh
join hosts h on vh.host_id = h.id
where vh.video_id = any(@video_ids::uuid[])
order by 1, vh.position;
