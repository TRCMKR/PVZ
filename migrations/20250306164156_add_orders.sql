-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders
(
    id              SERIAL PRIMARY KEY,
    user_id         INT              NOT NULL,
    weight          DOUBLE PRECISION NOT NULL,
    price           BIGINT           NOT NULL,
    packaging       int              NOT NULL,
    extra_packaging int              NOT NULL,
    status          VARCHAR(255),
    arrival_date    TIMESTAMP        NOT NULL,
    expiry_date     TIMESTAMP        NOT NULL,
    last_change     TIMESTAMP        NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE orders;
-- +goose StatementEnd
