-- name: GetResultUsingExecutionID :one
SELECT *
FROM Execution
WHERE execution_id = $1 LIMIT 1 ;

-- name: GetAllExecutions :many
SELECT * FROM Execution
ORDER BY executed_at;

-- name: InsertExecutionRequest :one
INSERT INTO Execution (
    execution_id,
    code,
    environment
) VALUES (
    $1, $2, $3
)
RETURNING execution_id;

-- name: InsertExecutionResult :one
UPDATE Execution
SET (
    result,
    execution_status,
    executed_at
) = ($1, $2, $3)
WHERE execution_id = $4
RETURNING execution_id;

-- name: GetAllExecutionResults :many
SELECT execution_id, result, execution_status, executed_at FROM Execution
ORDER BY executed_at;
