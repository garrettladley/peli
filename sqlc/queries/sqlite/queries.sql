-- name: GetAllRestaurants :many
SELECT * FROM restaurants ORDER BY position ASC;

-- name: GetRestaurantByID :one
SELECT * FROM restaurants WHERE id = ? LIMIT 1;

-- name: CreateRestaurant :one
INSERT INTO restaurants (name, cuisine, elo, score, position)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateRestaurant :exec
UPDATE restaurants
SET elo = ?, score = ?, position = ?
WHERE id = ?;

-- name: UpdatePositions :exec
UPDATE restaurants
SET position = position + 1
WHERE position >= ?;

-- name: RecordComparison :exec
INSERT INTO comparisons (winner_id, winner_name, loser_id, loser_name)
VALUES (?, ?, ?, ?);

-- name: CountRestaurants :one
SELECT COUNT(*) FROM restaurants;

-- name: DeleteRestaurant :exec
DELETE FROM restaurants WHERE id = ?;

-- name: GetAllComparisons :many
SELECT * FROM comparisons ORDER BY created_at DESC;

-- name: DeleteNonSeedRestaurants :exec
DELETE FROM restaurants WHERE is_seed = 0;

-- name: DeleteAllRestaurants :exec
DELETE FROM restaurants;

-- name: DeleteAllComparisons :exec
DELETE FROM comparisons;

-- name: CreateSeedRestaurant :one
INSERT INTO restaurants (name, cuisine, elo, score, position, is_seed)
VALUES (?, ?, ?, ?, ?, 1)
RETURNING *;
