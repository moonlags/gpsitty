-- name: GetUser :one
SELECT * FROM users where id = $1;

-- name: CreateUser :one
INSERT INTO users (id,name,email,avatar) VALUES ($1,$2,$3,$4) RETURNING *;

-- name: CreatePosition :exec
INSERT INTO positions (latitude,longitude,speed,heading,device_imei)
VALUES ($1,$2,$3,$4,$5);

-- name: UpdateBatteryPower :exec
UPDATE devices SET battery_power = $1 WHERE imei = $2;

-- name: InsertDevice :one
INSERT INTO devices (imei,battery_power,charging) VALUES ($1,$2,$3) RETURNING *;

-- name: LinkDevice :exec
INSERT INTO user_devices (userid,device_imei) VALUES ($1,$2);

-- name: UpdateCharging :exec
UPDATE devices SET charging = $1 where imei = $2;