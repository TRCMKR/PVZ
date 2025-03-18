-- +goose Up
-- +goose StatementBegin
INSERT INTO packagings(id, name)
VALUES
    (0, 'none'),
    (1, 'bax'),
    (2, 'box'),
    (3, 'wrap');

INSERT INTO statuses(id, name)
VALUES
    (1, 'stored'),
    (2, 'given'),
    (3, 'returned'),
    (4, 'deleted');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd
