-- +goose Up
alter table videos add youtube_id text,
add description_pl text,
add description_en text,
add constraint videos_youtube_id_must_be_set_for_yt_videos
check ((youtube_id is not null) = (provider = 'youtube'));

-- +goose Down
alter table videos drop youtube_id, drop description_pl, drop description_en;
