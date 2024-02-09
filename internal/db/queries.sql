-- name: GetFlags :many
SELECT * FROM flag ORDER BY name;

-- name: GetTrophies :many
SELECT * FROM trophy ORDER BY name;

-- name: GetClubs :many
SELECT * FROM entity WHERE type = 'CLUB' ORDER BY name;

-- name: GetEntities :many
SELECT * FROM entity ORDER BY name;

-- name: GetLeagues :many
SELECT * FROM league ORDER BY name;
