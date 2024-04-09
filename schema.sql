CREATE TABLE devices (
    imei VARCHAR(20) PRIMARY KEY,
    battery_power SMALLINT NOT NULL,
    charging BOOLEAN NOT NULL,
    last_status_packet TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE positions (
    id SERIAL PRIMARY KEY,
    latitude float8 NOT NULL,
    longitude float8 NOT NULL,
    speed SMALLINT NOT NULL,
    heading SMALLINT NOT NULL,
    device_imei VARCHAR(20) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (device_imei) REFERENCES devices(imei)
);
CREATE TABLE users (
    id VARCHAR(20) PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    email VARCHAR(50) UNIQUE NOT NULL,
    avatar VARCHAR(100) NOT NULL,
    last_login_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE user_devices (
    userid VARCHAR(20) NOT NULL,
    device_imei VARCHAR(20) NOT NULL,
    PRIMARY KEY (userid, device_imei),
    FOREIGN KEY (userid) REFERENCES users(id),
    FOREIGN KEY (device_imei) REFERENCES devices(imei)
);
