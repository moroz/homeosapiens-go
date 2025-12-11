-- +goose Up
-- +goose StatementBegin
truncate videos, video_sources;

alter table videos add column event_id uuid not null references events (id);
create index on videos (event_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table videos drop column event_id;
-- +goose StatementEnd
