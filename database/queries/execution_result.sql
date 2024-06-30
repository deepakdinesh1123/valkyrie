-- name: GetResultUsingExecutionID :one
SELECT *
FROM ExecutionResult
WHERE execution_id = $1 LIMIT 1 ;

-- name: GetAllExecutionResults :many
SELECT * FROM ExecutionResult
ORDER BY executed_at;

-- name: InsertExecutionResult :one
INSERT INTO ExecutionResult (
    execution_id,
    result,
    code,
    devenv,
    execution_status,
    executed_at
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;
