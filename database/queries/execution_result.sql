-- name: GetResultUsingSandboxID :one
SELECT result
FROM ExecutionResult
WHERE sandbox = $1 LIMIT 1 ;

-- name: GetAllExecutionResults :many
SELECT * FROM ExecutionResult
ORDER BY executed_at;

-- name: InsertExecutionResult :one
INSERT INTO ExecutionResult (
    result,
    cmdLineArgs,
    environment,
    flake,
    dependency_file, 
    sandbox
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;
