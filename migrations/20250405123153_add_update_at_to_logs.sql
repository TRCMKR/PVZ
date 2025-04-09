-- +goose Up
-- +goose StatementBegin
ALTER TABLE logs
    ADD COLUMN updated_at timestamp DEFAULT now();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE logs
     DROP COLUMN updated_at;
-- +goose StatementEnd
