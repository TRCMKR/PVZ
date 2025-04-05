-- +goose Up
-- +goose StatementBegin
INSERT INTO job_statuses(name)
VALUES
    ('CREATED'),
    ('PROCESSING'),
    ('FAILED'),
    ('NO_ATTEMPTS_LEFT');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM job_statuses
WHERE name IN ('CREATED', 'PROCESSING', 'FAILED', 'NO_ATTEMPTS_LEFT');
-- +goose StatementEnd
