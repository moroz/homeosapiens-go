-- +goose Up
-- +goose StatementBegin
alter table users add column profile_picture text;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table users drop column profile_picture;
-- +goose StatementEnd
