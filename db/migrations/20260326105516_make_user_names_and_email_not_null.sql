-- +goose Up
-- +goose StatementBegin
delete from users where given_name_encrypted is null or family_name_encrypted is null or email_encrypted is null;

alter table users
    alter column given_name_encrypted set not null,
    alter column family_name_encrypted set not null,
    alter column email_encrypted set not null;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table users
    alter column given_name_encrypted drop not null,
    alter column family_name_encrypted drop not null,
    alter column email_encrypted drop not null;
-- +goose StatementEnd
