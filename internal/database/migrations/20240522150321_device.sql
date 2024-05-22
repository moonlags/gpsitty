-- +goose Up
-- +goose StatementBegin
CREATE TABLE devices (
    imei VARCHAR(15) PRIMARY KEY,
    battery_power SMALLINT NOT NULL,
    charging BOOLEAN NOT NULL,
    last_status_packet TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE devices;
-- +goose StatementEnd
