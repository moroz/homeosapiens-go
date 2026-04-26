-- +goose Up
alter table videos
    add recorded_on date,
    add host_id uuid references hosts(id) on delete set null,
    drop thumbnail_id,
    add thumbnail_en_id uuid references assets (id) on delete set null,
    add thumbnail_pl_id uuid references assets (id) on delete set null;

-- +goose Down
alter table videos
    drop recorded_on,
    drop host_id,
    drop thumbnail_en_id,
    drop thumbnail_pl_id,
    add thumbnail_id uuid references assets (id) on delete set null;
