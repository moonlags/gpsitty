-- +goose Up
-- +goose StatementBegin
CREATE TABLE positions (
    id SERIAL PRIMARY KEY,
    latitude float8 NOT NULL,
    longitude float8 NOT NULL,
    speed SMALLINT NOT NULL,
    heading SMALLINT NOT NULL,
    device_imei VARCHAR(15) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (device_imei) REFERENCES devices(imei)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE positions;
-- +goose StatementEnd
