-- +goose Up
-- +goose StatementBegin
alter table hosts add column country char(2);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table hosts drop column country;
-- +goose StatementEnd
