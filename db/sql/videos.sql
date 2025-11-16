-- name: ListVideos :many
select * from videos v
order by v.id desc;

-- name: ListVideoSourcesForVideos :many
select * from video_sources vs
where vs.video_id = any(@VideoIDs::uuid[])
order by vs.video_id, vs.id;
