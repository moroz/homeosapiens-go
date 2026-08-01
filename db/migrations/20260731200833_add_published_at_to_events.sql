-- +goose Up
alter table events
add published_at timestamp,
alter column description_en drop not null,
alter column description_pl drop not null,
add constraint published_event_must_have_description check ((published_at is null) OR (description_en is not null and description_pl is not null));


-- +goose Down
alter table events
drop published_at,
alter column description_en set not null,
alter column description_pl set not null;
