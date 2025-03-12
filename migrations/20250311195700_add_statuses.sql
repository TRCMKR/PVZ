-- +goose Up
-- +goose StatementBegin
CREATE TABLE statuses
(
    id   serial PRIMARY KEY,
    name varchar(100)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE statuses;
-- +goose StatementEnd
