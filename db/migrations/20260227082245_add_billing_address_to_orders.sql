-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders
    ADD COLUMN billing_given_name_encrypted       BYTEA NOT NULL,
    ADD COLUMN billing_family_name_encrypted        BYTEA NOT NULL,
    ADD COLUMN billing_phone_encrypted            BYTEA,
    ADD COLUMN billing_street_encrypted           BYTEA NOT NULL,
    ADD COLUMN billing_house_number_encrypted     BYTEA NOT NULL,
    ADD COLUMN billing_apartment_number_encrypted BYTEA,
    ADD COLUMN billing_city_encrypted             BYTEA NOT NULL,
    ADD COLUMN billing_postal_code_encrypted      BYTEA,
    ADD COLUMN billing_country                    CHAR(2) NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE orders
    DROP COLUMN billing_given_name_encrypted,
    DROP COLUMN billing_family_name_encrypted,
    DROP COLUMN billing_phone_encrypted,
    DROP COLUMN billing_street_encrypted,
    DROP COLUMN billing_house_number_encrypted,
    DROP COLUMN billing_apartment_number_encrypted,
    DROP COLUMN billing_city_encrypted,
    DROP COLUMN billing_postal_code_encrypted,
    DROP COLUMN billing_country;
-- +goose StatementEnd
