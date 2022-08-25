-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN superuser BOOLEAN NOT NULL DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN superuser;
-- +goose StatementEnd
