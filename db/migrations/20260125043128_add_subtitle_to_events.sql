-- +goose Up
-- +goose StatementBegin
alter table events
add column subtitle_en varchar(255),
add column subtitle_pl varchar(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table events drop column subtitle_en, drop column subtitle_pl;
-- +goose StatementEnd
