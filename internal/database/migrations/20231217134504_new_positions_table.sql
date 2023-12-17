-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS positions(
    latitude REAL,
    longitude REAL,
    speed smallint,
    heading integer,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS positions;
-- +goose StatementEnd
