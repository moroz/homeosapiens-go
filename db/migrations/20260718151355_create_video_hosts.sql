-- +goose Up
create table video_hosts (
  id uuid not null primary key default (uuidv7()),
  video_id uuid not null,
  host_id uuid not null,
  position int not null,
  inserted_at timestamp not null default (now()),
  updated_at timestamp not null default (now()),
  check (position > 0),
  unique (video_id, host_id),
  unique (video_id, position)
);

insert into video_hosts (video_id, host_id, position)
select id, host_id, 1 from videos where host_id is not null;

-- +goose Down
drop table video_hosts;
