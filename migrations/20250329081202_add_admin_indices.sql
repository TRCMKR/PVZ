-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_admin_id_hash ON admins USING HASH (id);
CREATE INDEX idx_admin_username_hash ON admins USING HASH (username);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_admin_id_hash;
DROP INDEX IF EXISTS idx_admin_username_hash;
-- +goose StatementEnd
