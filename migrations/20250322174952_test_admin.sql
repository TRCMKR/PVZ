-- +goose Up
-- +goose StatementBegin
INSERT INTO admins(username, password)
VALUES ('test', '$2a$10$XRHPGT8WCS5NLyZpvk0dUO/9u85Ng7oeuZVNnzShj6SCXWCCirHV6')
-- password 12345678
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
