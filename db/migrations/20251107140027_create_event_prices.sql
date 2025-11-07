-- +goose Up
-- +goose StatementBegin
CREATE TYPE price_type AS ENUM ('Fixed', 'Percentage', 'DiscountFixed', 'DiscountPercentage');
CREATE TYPE price_rule_type AS ENUM ('Base', 'DiscountCode', 'EarlyBird');
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE event_prices (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    price_type price_type NOT NULL DEFAULT 'Fixed',
    rule_type price_rule_type NOT NULL DEFAULT 'Base',
    price_amount DECIMAL(20, 8) NOT NULL,
    price_currency VARCHAR(3) NOT NULL DEFAULT 'PLN',
    discount_code CITEXT,
    priority INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    valid_from TIMESTAMP(0),
    valid_until TIMESTAMP(0),
    created_at TIMESTAMP(0) NOT NULL DEFAULT (now() at time zone 'utc'),
    updated_at TIMESTAMP(0) NOT NULL DEFAULT (now() at time zone 'utc'),
    check ((valid_until is null) = (rule_type != 'EarlyBird')),
    check ((discount_code is null) = (rule_type != 'DiscountCode'))
);

CREATE UNIQUE INDEX ON event_prices(discount_code, event_id);
CREATE INDEX ON event_prices(event_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS event_prices;
DROP TYPE IF EXISTS price_type;
DROP TYPE IF EXISTS price_rule_type;
-- +goose StatementEnd
