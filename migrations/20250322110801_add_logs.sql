-- +goose Up
-- +goose StatementBegin
CREATE TABLE logs
(
    id          SERIAL PRIMARY KEY,
    order_id    INT NOT NULL,
    admin_id    INT NOT NULL,
    message     TEXT,
    date        TIMESTAMP,
    url         TEXT,
    method      TEXT,
    status      INT,

    CONSTRAINT fk_logs_admin_id FOREIGN KEY (admin_id) REFERENCES admins (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE logs;
-- +goose StatementEnd
