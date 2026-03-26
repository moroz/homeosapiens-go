-- +goose Up
-- +goose StatementBegin
delete from orders;

alter table orders
    drop column billing_street_encrypted,
    drop column billing_house_number_encrypted,
    drop column billing_apartment_number_encrypted,
    add column billing_address_line1_encrypted bytea not null,
    add column billing_address_line2_encrypted bytea;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
delete from orders;

alter table orders
    drop column billing_address_line1_encrypted,
    drop column billing_address_line2_encrypted,
    add column billing_street_encrypted bytea not null,
    add column billing_house_number_encrypted bytea not null,
    add column billing_apartment_number_encrypted bytea;
-- +goose StatementEnd
