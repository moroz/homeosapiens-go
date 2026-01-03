-- +goose Up
-- +goose StatementBegin
alter table event_registrations
  add column given_name_encrypted bytea not null,
  add column family_name_encrypted bytea not null,
  add column email_encrypted bytea not null,
  add column country char(2) not null;

create unique index on event_registrations (event_id, user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table event_registrations
  drop column given_name_encrypted,
  drop column family_name_encrypted,
  drop column email_encrypted,
  drop column country;

drop index event_registrations_event_id_user_id_idx;
-- +goose StatementEnd
