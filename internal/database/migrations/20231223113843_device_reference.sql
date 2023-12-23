-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS positions;
CREATE TABLE IF NOT EXISTS positions(
    latitude REAL NOT NULL,
    longitude REAL NOT NULL,
    speed SMALLINT NOT NULL,
    heading INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    device_imei VARCHAR(255) NOT NULL REFERENCES devices(imei)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS positions;
CREATE TABLE IF NOT EXISTS positions(
    latitude REAL NOT NULL,
    longitude REAL NOT NULL,
    speed SMALLINT NOT NULL,
    heading INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd
