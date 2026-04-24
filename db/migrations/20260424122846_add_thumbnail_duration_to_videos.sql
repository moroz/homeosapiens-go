-- +goose Up
alter table videos add duration_seconds integer check (duration_seconds >= 0),
    add thumbnail_id uuid references assets (id) on delete set null;

-- +goose Down
alter table videos drop duration_seconds, drop thumbnail_id;
