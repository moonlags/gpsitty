-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS devices (
    imei TEXT PRIMARY KEY,
    battery_power INTEGER NOT NULL,
    charging INTEGER NOT NULL,
    last_status_packet TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS positions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    latitude REAL NOT NULL,
    longitude REAL NOT NULL,
    speed INTEGER NOT NULL,
    heading INTEGER NOT NULL,
    device_imei TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (device_imei) REFERENCES devices(imei)
);
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE,
    avatar TEXT NOT NULL,
    last_login_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS user_devices (
    userid TEXT NOT NULL,
    device_imei TEXT NOT NULL,
    PRIMARY KEY (userid, device_imei),
    FOREIGN KEY (userid) REFERENCES users(id),
    FOREIGN KEY (device_imei) REFERENCES devices(imei)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_devices;
DROP TABLE IF EXISTS positions;
DROP TABLE IF EXISTS devices;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
