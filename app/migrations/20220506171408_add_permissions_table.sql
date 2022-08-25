-- +goose Up
-- +goose StatementBegin
CREATE TABLE permissions (
    name text PRIMARY KEY
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE permissions;
-- +goose StatementEnd
