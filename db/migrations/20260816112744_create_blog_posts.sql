-- +goose Up
create table blog_posts (
  id uuid not null default uuidv7(),
  title text not null,
  slug citext not null unique,
  language locale not null default 'en',
  body text not null,
  published_at timestamp(0),
  inserted_at timestamp(0) not null default now(),
  updated_at timestamp(0) not null default now()
);

-- +goose Down
drop table blog_posts;
