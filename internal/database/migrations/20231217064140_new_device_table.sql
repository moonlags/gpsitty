-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS devices(
    imei VARCHAR(255) PRIMARY KEY,
    userid VARCHAR(255) NOT NULL,
    batery_power smallint,
    last_status_packet TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    status_cooldown integer,
    logged_in boolean,
    charging boolean
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS devices;
-- +goose StatementEnd
