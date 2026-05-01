-- +goose Up
-- +goose StatementBegin
alter table assets
add scaled boolean not null default false,
alter column object_key drop not null,
add constraint assets_object_key_required_for_regular_assets check (
    (object_key is null) = scaled
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
delete from assets where object_key is null;
alter table assets drop scaled, alter column object_key set not null;
-- +goose StatementEnd
