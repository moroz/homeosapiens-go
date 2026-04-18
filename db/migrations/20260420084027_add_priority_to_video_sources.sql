-- +goose Up
-- +goose StatementBegin
truncate video_sources;
alter table video_sources add priority int not null;
create unique index on video_sources (video_id, priority);
alter table video_sources add constraint video_sources_priority_must_be_non_neg check (priority >= 0);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table video_sources drop priority;
-- +goose StatementEnd
