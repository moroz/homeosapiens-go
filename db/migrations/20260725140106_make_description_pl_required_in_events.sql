-- +goose Up
alter table events alter column description_pl set not null;

-- +goose Down
alter table events alter column description_pl drop not null;
