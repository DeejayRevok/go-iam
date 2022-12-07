-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id VARCHAR(36) PRIMARY KEY,
    username text UNIQUE,
    email text UNIQUE,
    password text
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
