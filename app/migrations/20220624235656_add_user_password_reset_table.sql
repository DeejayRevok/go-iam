-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_password_resets (
    token text PRIMARY KEY,
    expiration timestamp NOT NULL ,
    user_id VARCHAR(36) UNIQUE NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_password_resets;
-- +goose StatementEnd
