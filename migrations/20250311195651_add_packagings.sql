-- +goose Up
-- +goose StatementBegin
CREATE TABLE packagings
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR(100)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE packagings;
-- +goose StatementEnd
