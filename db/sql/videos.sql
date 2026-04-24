-- name: ListVideos :many
select v.* from videos v
order by v.id desc;

-- name: ListVideoGroupsForUser :many
select vg.* from video_groups vg
where vg.product_id is null
or vg.id in (
  select vg.id from video_groups vg
  join user_product_access upa on upa.product_id = vg.product_id
  where upa.user_id = sqlc.narg(user_id)::uuid
) or exists (
  select from users u where u.id = sqlc.narg(user_id)::uuid and u.user_role = 'Administrator'
);

-- name: GetVideoGroupForUserBySlug :one
select vg.* from video_groups vg
where (sqlc.narg(slug)::text is null or vg.slug = sqlc.narg(slug)) and (
    vg.product_id is null or vg.id in (
        select vg.id from video_groups vg
        join user_product_access upa on upa.product_id = vg.product_id
        where upa.user_id = sqlc.narg(user_id)::uuid
    ) or exists (
        select from users u where u.id = sqlc.narg(user_id)::uuid and u.user_role = 'Administrator'
    )
) limit 1;

-- name: GetVideoForUser :one
select v.* from videos v
join video_groups_videos vgv on vgv.video_id = v.id
join video_groups vg on vgv.video_group_id = vg.id
where v.slug = @video_slug and vg.slug = @group_slug and (
    vg.product_id is null or vg.id in (
        select vg.id from video_groups vg
        join user_product_access upa on upa.product_id = vg.product_id
        where upa.user_id = sqlc.narg(user_id)::uuid
    ) or exists (
        select from users u where u.id = sqlc.narg(user_id)::uuid and u.user_role = 'Administrator'
    )
) limit 1;

-- name: ListVideosForVideoGroup :many
select sqlc.embed(v), sqlc.embed(a) from videos v
join video_groups_videos vgv on vgv.video_id = v.id
left join assets a on v.thumbnail_id = a.id
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
