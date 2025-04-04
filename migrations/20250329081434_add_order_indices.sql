-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_orders_id_hash ON orders USING HASH (id);
CREATE INDEX idx_orders_user_id_hash ON orders USING HASH (user_id);
CREATE INDEX idx_orders_dates ON orders (arrival_date, expiry_date);
CREATE INDEX idx_orders_last_change ON orders (last_change);
CREATE INDEX idx_orders_status ON orders (status);
CREATE INDEX idx_orders_status_hash ON orders USING HASH (status);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_orders_id_hash;
DROP INDEX IF EXISTS idx_orders_user_id_hash;
DROP INDEX IF EXISTS idx_orders_dates;
DROP INDEX IF EXISTS idx_orders_last_change;
DROP INDEX IF EXISTS idx_orders_status;
DROP INDEX IF EXISTS idx_orders_status_hash;
-- +goose StatementEnd
