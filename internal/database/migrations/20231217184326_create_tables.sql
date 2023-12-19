-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS devices(
    imei VARCHAR(255) PRIMARY KEY,
    battery_power SMALLINT NOT NULL,
    last_status_packet TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    status_cooldown INTEGER NOT NULL,
    charging BOOLEAN NOT NULL
);
CREATE TABLE IF NOT EXISTS positions(
    latitude REAL NOT NULL,
    longitude REAL NOT NULL,
    speed SMALLINT NOT NULL,
    heading INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS users(
    id VARCHAR(255) PRIMARY KEY,
    email VARCHAR(255) UNIQUE,
    avatar VARCHAR(255) NOT NULL
);
CREATE TABLE IF NOT EXISTS user_devices(
    userid VARCHAR(255) PRIMARY KEY REFERENCES users(id),
    device_imei VARCHAR(255) NOT NULL REFERENCES devices(imei)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_devices;
DROP TABLE IF EXISTS positions;
DROP TABLE IF EXISTS devices;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
