-- +goose Up
-- +goose StatementBegin
INSERT INTO job_statuses(name)
VALUES
    ('DONE');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM job_statuses
WHERE name = 'DONE';
-- +goose StatementEnd
