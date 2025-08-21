-- name: UpsertUser :one
INSERT INTO users (username, email)
VALUES ($1, $2)
ON CONFLICT (username) DO UPDATE SET
    email = EXCLUDED.email,
    modified_at = NOW()
RETURNING id, username, email, created_at, modified_at;

-- name: GetUser :one
SELECT id, username, email, created_at, modified_at
FROM users
WHERE
  id = sqlc.narg('id') OR
  email = sqlc.narg('email') OR
  username = sqlc.narg('username')
LIMIT 1;
