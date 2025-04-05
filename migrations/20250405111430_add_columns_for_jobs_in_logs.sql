-- +goose Up
-- +goose StatementBegin
ALTER TABLE logs
    ADD COLUMN job_status INT DEFAULT 1,
    ADD COLUMN attempts_left INT DEFAULT 3,
    ADD CONSTRAINT fk_job_status FOREIGN KEY (job_status) REFERENCES job_statuses(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE logs
    DROP CONSTRAINT fk_job_status,
    DROP COLUMN job_status,
    DROP COLUMN attempts_left;
-- +goose StatementEnd
