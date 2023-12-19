-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS devices(
    imei VARCHAR(255) PRIMARY KEY,
    battery_power SMALLINT,
    last_status_packet TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    status_cooldown INTEGER,
    charging BOOLEAN
);
CREATE TABLE IF NOT EXISTS positions(
    latitude REAL,
    longitude REAL,
    speed SMALLINT,
    heading INTEGER,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS user_devices(
    userid VARCHAR(255) PRIMARY KEY,
    device_imei VARCHAR(255) NOT NULL,
    FOREIGN KEY (userid) REFERENCES users(id),
    FOREIGN KEY (device_imei) REFERENCES devices(imei)
);
CREATE TABLE IF NOT EXISTS users(
    id VARCHAR(255) PRIMARY KEY,
    email VARCHAR(255) UNIQUE,
    avatar VARCHAR(255)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS positions;
DROP TABLE IF EXISTS devices;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS user_devices;
-- +goose StatementEnd
