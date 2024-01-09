-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS devices (
    imei VARCHAR(20) PRIMARY KEY,
    battery_power INT NOT NULL,
    status_cooldown INT NOT NULL,
    charging BOOLEAN NOT NULL,
    last_status_packet TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS positions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    latitude DOUBLE NOT NULL,
    longitude DOUBLE NOT NULL,
    speed SMALLINT NOT NULL,
    heading INTEGER NOT NULL,
    device_imei VARCHAR(20) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (device_imei) REFERENCES devices(imei)
);
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(255) PRIMARY KEY,
    email VARCHAR(255) UNIQUE,
    avatar VARCHAR(255) NOT NULL,
    last_login_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS user_devices (
    userid VARCHAR(255) NOT NULL,
    device_imei VARCHAR(20) NOT NULL,
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
