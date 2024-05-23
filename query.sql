-- name: GetUser :one
SELECT * FROM users WHERE id = ?;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = ?;

-- name: CreateUser :exec
INSERT INTO users (id,email,password) VALUES (?,?,?);

-- name: CreatePosition :exec
INSERT INTO positions (latitude,longitude,speed,heading,device_imei) VALUES (?,?,?,?,?);

-- name: UpdateBatteryPower :exec
UPDATE devices SET battery_power = ? WHERE imei = ?;

-- name: InsertDevice :exec
INSERT INTO devices (imei,battery_power,charging) VALUES (?,?,?) ON CONFLICT DO NOTHING;

-- name: LinkDevice :exec
INSERT INTO user_devices (userid,device_imei) VALUES (?,?);

-- name: UpdateCharging :exec
UPDATE devices SET charging = ? WHERE imei = ?;

-- name: GetDevices :many
SELECT devices.* from devices JOIN user_devices ON devices.imei = user_devices.device_imei WHERE user_devices.userid = ?;