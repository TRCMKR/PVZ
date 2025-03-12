-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders
    ADD CONSTRAINT fk_orders_packaging FOREIGN KEY (packaging) REFERENCES packagings (id) ON DELETE CASCADE,
    ADD CONSTRAINT fk_orders_extra_packaging FOREIGN KEY (extra_packaging) REFERENCES packagings (id) ON DELETE CASCADE,
    ADD CONSTRAINT fk_orders_status FOREIGN KEY (status) REFERENCES statuses (id) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE orders
    DROP CONSTRAINT fk_orders_packaging,
    DROP CONSTRAINT fk_orders_extra_packaging,
    DROP CONSTRAINT fk_orders_status;
-- +goose StatementEnd
