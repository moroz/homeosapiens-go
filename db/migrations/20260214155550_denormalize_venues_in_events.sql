-- +goose Up
-- +goose StatementBegin
alter table events
    add column venue_name_en varchar(255),
    add column venue_name_pl varchar(255),
    add column venue_street varchar(255),
    add column venue_city_en varchar(255),
    add column venue_city_pl varchar(255),
    add column venue_postal_code varchar(15),
    add column venue_country_code varchar(2);

update events e
set
    venue_name_en = v.name_en,
    venue_name_pl = v.name_pl,
    venue_street = v.street,
    venue_city_en = v.city_en,
    venue_city_pl = v.city_pl,
    venue_postal_code = v.postal_code,
    venue_country_code = v.country_code
from venues v
where e.venue_id = v.id;

alter table events
    add constraint events_venue_required_for_non_virtual_events check (
        is_virtual = true or (
            venue_name_en is not null and
            venue_street is not null and
            venue_city_en is not null and
            venue_country_code is not null
            )
        );

alter table events drop column venue_id;

drop table venues;
-- +goose StatementEnd
