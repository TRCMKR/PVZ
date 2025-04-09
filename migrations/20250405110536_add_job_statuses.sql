-- +goose Up
-- +goose StatementBegin
CREATE TABLE job_statuses(
    id SERIAL PRIMARY KEY,
    name TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE job_statuses
-- +goose StatementEnd
