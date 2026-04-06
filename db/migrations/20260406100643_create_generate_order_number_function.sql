-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION generate_order_number(for_date date DEFAULT (NOW() AT TIME ZONE 'Europe/Warsaw')::date)
RETURNS bigint
LANGUAGE plpgsql
AS $$
DECLARE
    v_prefix bigint;
    v_next   bigint;
BEGIN
    v_prefix :=
         (EXTRACT(YEAR  FROM for_date)::bigint % 100) * 100000000  -- YY
        + EXTRACT(MONTH FROM for_date)::bigint        *   1000000  -- MM
        + EXTRACT(DAY   FROM for_date)::bigint        *     10000; -- DD

    -- Obtains an exclusive transaction-level advisory lock. This lock gets implicitly unlocked when a transaction
    -- is committed or cancelled. All concurrent transactions calling the same function will wait until this lock
    -- gets released, which also means that only one order can be processed at a time.
    PERFORM pg_advisory_xact_lock(v_prefix);

    SELECT COALESCE(MAX(order_number % 10000), 0) + 1 + v_prefix
    INTO v_next
    FROM orders
    WHERE order_number BETWEEN v_prefix AND v_prefix + 9999;

    RETURN v_next;
END;
$$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP FUNCTION IF EXISTS generate_order_number(date);
-- +goose StatementEnd
