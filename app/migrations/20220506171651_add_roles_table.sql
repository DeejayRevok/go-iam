-- +goose Up
-- +goose StatementBegin
CREATE TABLE roles (
    id VARCHAR(36) PRIMARY KEY,
    name text UNIQUE
);
CREATE TABLE role_permission (
    role_id VARCHAR(36) NOT NULL,
    permission_name TEXT NOT NULL,
    FOREIGN KEY (role_id) REFERENCES roles(id), 
    FOREIGN KEY (permission_name) REFERENCES permissions(name),
    UNIQUE (role_id, permission_name)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE role_permission;
DROP TABLE roles;
-- +goose StatementEnd
