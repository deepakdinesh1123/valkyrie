CREATE TABLE IF NOT EXISTS Worker(
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(50),
    modified_at TIMESTAMP,
    modified_by VARCHAR(50)
);

-- name: InsertWorker :one
INSERT INTO Worker
    (name)
VALUES
    ($1)
RETURNING *;

-- name: GetWorker :one
SELECT * FROM Worker
WHERE name = $1
LIMIT 1;

-- name: GetAllWorkers :many
SELECT * FROM Worker
LIMIT $1 OFFSET $2;

-- name: GetTotalWorkers :one
SELECT count(*) FROM Worker;