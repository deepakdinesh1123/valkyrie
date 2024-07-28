CREATE TABLE IF NOT EXISTS JobQueue(
    id bigserial PRIMARY KEY,
    created_by VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    script TEXT,
    script_path VARCHAR(50),
    args jsonb,
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
    (script, flake, language, script_path)
VALUES
    ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateJob :one
UPDATE JobQueue
SET
    completed_at = current_timestamp
    , logs = CONCAT(logs, $2)
WHERE id = $1 AND completed_at IS NULL
RETURNING *;

-- name: RemainingJobs :one
SELECT count(*) FROM JobQueue
WHERE queue=$1 AND completed_at IS NULL;


-- name: GetAllJobs :many
SELECT * FROM JobQueue
ORDER BY started_at;

-- name: GetJob :one
SELECT * FROM JobQueue
WHERE id = $1
LIMIT 1;

-- name: GetAllExecutions :many
SELECT * FROM JobQueue
WHERE job_type = 'execution'
ORDER BY started_at;

-- name: GetResultUsingExecutionID :one
SELECT *
FROM JobQueue
WHERE id = $1 AND job_type = 'execution' LIMIT 1 ;

-- name: GetAllExecutionResults :many
SELECT id, logs, script, args FROM JobQueue
WHERE job_type = 'execution'
ORDER BY started_at;