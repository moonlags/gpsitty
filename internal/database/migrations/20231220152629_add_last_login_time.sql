-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN last_login_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN last_login_time;
-- +goose StatementEnd
