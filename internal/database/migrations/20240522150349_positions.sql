-- +goose Up
-- +goose StatementBegin
CREATE TABLE positions (
    id SERIAL PRIMARY KEY,
    latitude REAL NOT NULL,
    longitude REAL NOT NULL,
    speed INT NOT NULL,
    heading INT NOT NULL,
    device_imei VARCHAR(15) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (device_imei) REFERENCES devices(imei)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE positions;
-- +goose StatementEnd
