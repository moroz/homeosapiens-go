-- +goose Up
-- +goose StatementBegin
create type user_role as enum ('Regular', 'Administrator');
alter table users add column user_role user_role not null default 'Regular';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table users drop column user_role;
drop type user_role;
-- +goose StatementEnd
