-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id VARCHAR(36) PRIMARY KEY,
    username text UNIQUE,
    email text UNIQUE,
    password text
);
CREATE TABLE user_role (
    user_id VARCHAR(36) NOT NULL,
    role_id VARCHAR(36) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (role_id) REFERENCES roles(id), 
    UNIQUE (user_id, role_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_role;
DROP TABLE users;
-- +goose StatementEnd
