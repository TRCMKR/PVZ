-- +goose Up
-- +goose StatementBegin
DO $$
    DECLARE i INT := 0;
    BEGIN
        WHILE i < 100000 LOOP
            INSERT INTO orders (user_id, weight, price, packaging, extra_packaging, status, arrival_date, expiry_date, last_change)
            VALUES (
               floor(random() * 1000)::int,
               random() * 100,
               floor(random() * 10000)::bigint,
               0,
               0,
               1,
               now(),
               now() + interval '100 days',
               now()
                    )
            ON CONFLICT (id) DO NOTHING;
            i := i + 1;
    END LOOP;
END $$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- DELETE FROM orders
-- WHERE id IN (
--     SELECT id FROM orders
--     ORDER BY id DESC
--     LIMIT 100000
-- );
-- ALTER SEQUENCE orders_id_seq RESTART WITH 5;
-- +goose StatementEnd
