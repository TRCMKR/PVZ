-- +goose Up
-- +goose StatementBegin
INSERT INTO orders(user_id, weight, price, packaging, extra_packaging, status, arrival_date, expiry_date, last_change)
    VALUES
        (52, 100, 1233, 2, 3, 1, '2022-03-20', '2024-04-01', '2022-03-20 12:30:00'),
        (789, 1233, 22222, 3, 0, 2, '2022-03-20', '7022-04-01', '2025-03-17 12:30:00'),
        (22, 1233, 22222, 0, 0, 3, '2022-03-20', '7022-04-01', '2022-03-20 12:30:00'),
        (789, 1233, 22222, 0, 0, 1, '2022-03-20', '7022-04-01', '2022-03-20 12:30:00')
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd
