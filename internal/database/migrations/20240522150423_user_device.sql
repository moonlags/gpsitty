-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_devices (
    userid TEXT NOT NULL,
    device_imei VARCHAR(15) NOT NULL,
    PRIMARY KEY (userid, device_imei),
    FOREIGN KEY (userid) REFERENCES users(id),
    FOREIGN KEY (device_imei) REFERENCES devices(imei)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_devices;
-- +goose StatementEnd
