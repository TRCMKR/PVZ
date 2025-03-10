-- +goose Up
-- +goose StatementBegin
create table orders(
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    weight DOUBLE PRECISION not null,
    price BIGINT not null ,
    packaging int not null,
    extra_packaging int not null,
    status VARCHAR(255),
    arrival_date TIMESTAMP not null,
    expiry_date TIMESTAMP not null,
    last_change TIMESTAMP not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE orders;
-- +goose StatementEnd
