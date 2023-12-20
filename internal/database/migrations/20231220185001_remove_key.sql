-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS user_devices;
CREATE TABLE IF NOT EXISTS user_devices(
    userid VARCHAR(255) REFERENCES users(id),
    device_imei VARCHAR(255) NOT NULL REFERENCES devices(imei)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_devices;
CREATE TABLE IF NOT EXISTS user_devices(
    userid VARCHAR(255) PRIMARY KEY REFERENCES users(id),
    device_imei VARCHAR(255) NOT NULL REFERENCES devices(imei)
);
-- +goose StatementEnd
