CREATE TABLE IF NOT EXISTS JobQueue(
    id bigserial PRIMARY KEY,
    created_by VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    script TEXT,
    script_path VARCHAR(50),
    args VARCHAR(100),
    logs TEXT,
    flake TEXT,
    language VARCHAR(50),
    mem_peak INT,
    timeout INT,
    priority INT,
    lease_timeout FLOAT,
    queue VARCHAR(50) NOT NULL DEFAULT 'default',
    job_type VARCHAR(50) NOT NULL DEFAULT 'execution',
    worker_id INT REFERENCES Worker(id)
);

-- name: FetchJob :one
UPDATE JobQueue SET started_at = current_timestamp, worker_id = $1
WHERE id = (
    SELECT id FROM JobQueue
    WHERE (completed_at IS NULL and started_at IS NULL)
    ORDER BY
        priority ASC,
        id ASC
    FOR UPDATE SKIP LOCKED
    LIMIT 1
    )
RETURNING *;

-- name: InsertJob :one
INSERT INTO JobQueue
    (script, flake, language, script_path, args)
VALUES
    ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateJob :one
UPDATE JobQueue
SET
    completed_at = current_timestamp
    , logs = $2
WHERE id = $1 AND completed_at IS NULL
RETURNING *;

-- name: GetAllJobs :many
SELECT *, count(*) OVER() AS total FROM JobQueue
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetJob :one
SELECT * FROM JobQueue
WHERE id = $1
LIMIT 1;

-- name: GetResultUsingExecutionID :one
SELECT *
FROM JobQueue
WHERE id = $1 AND job_type = 'execution' LIMIT 1 ;

-- name: DeleteJob :exec
DELETE FROM JobQueue
WHERE id = $1 and completed_at IS NULL;

-- name: GetTotalJobs :one
SELECT count(*) FROM JobQueue;