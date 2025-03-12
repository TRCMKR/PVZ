-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders
    ALTER COLUMN status TYPE INT USING status::INTEGER;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE orders
    ALTER COLUMN status TYPE VARCHAR(255) USING status::TEXT;
-- +goose StatementEnd
