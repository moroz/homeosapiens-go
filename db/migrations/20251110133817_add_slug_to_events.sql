-- +goose Up
-- +goose StatementBegin
truncate events cascade;
alter table events add column slug citext not null;
create unique index on events (slug);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table events drop column slug;
-- +goose StatementEnd
