-- name: CreateEnvironment :one
INSERT INTO environment (name) VALUES (?)
RETURNING *;

-- name: GetEnvironment :one
SELECT * FROM environment WHERE id = ? LIMIT 1;

-- name: GetAllEnvironments :many
SELECT * FROM environment;

-- name: DeleteEnvironment :exec
DELETE FROM environment WHERE id = ?;

-- name: GetEnvironmentByName :one
SELECT * FROM environment WHERE name = ? LIMIT 1;

-- name: CreateValue :one
INSERT INTO environment_values (environment_id, key, value) VALUES (?, ?, ?)
RETURNING *;

-- name: GetValue :one
SELECT * FROM environment_values WHERE id = ? LIMIT 1;

-- name: GetAllValues :many
SELECT * FROM environment_values;

-- name: DeleteValue :exec
DELETE FROM environment_values WHERE id = ?;

-- name: GetValuesByEnvironmentID :many
SELECT * FROM environment_values WHERE environment_id = ?;

-- name: GetValueByKey :one
SELECT * FROM environment_values WHERE environment_id = ? AND key = ? LIMIT 1;

-- name: UpdateValue :one
UPDATE environment_values SET value = ? WHERE id = ? RETURNING *;