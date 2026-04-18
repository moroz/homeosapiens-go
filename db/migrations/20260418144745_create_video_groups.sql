-- +goose Up
-- +goose StatementBegin
create table video_groups (
  id uuid not null primary key default uuidv7(),
  title_en text not null,
  title_pl text not null,
  slug citext not null unique,
  product_id uuid references products (id) on delete set null,
  inserted_at timestamp not null default now(),
  updated_at timestamp not null default now()
);

create index on video_groups (product_id);

create table video_groups_videos (
  id uuid not null primary key default uuidv7(),
  position integer not null,
  video_id uuid not null references videos (id) on delete cascade,
  video_group_id uuid not null references video_groups (id) on delete cascade,
  inserted_at timestamp not null default now(),
  updated_at timestamp not null default now(),
  unique (video_id, video_group_id),
  unique (video_group_id, position),
  check (position >= 0)
);

alter table videos drop column event_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
truncate videos cascade;
alter table videos add column event_id uuid not null references events (id) on delete cascade;

drop table video_groups_videos;
drop table video_groups;
-- +goose StatementEnd
