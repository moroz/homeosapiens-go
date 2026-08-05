-- +goose Up
alter table events add meeting_url text;

-- +goose Down
alter table events drop meeting_url;
