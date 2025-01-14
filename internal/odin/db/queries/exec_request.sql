-- name: InsertExecRequest :one
insert into exec_request
    (
        hash, 
        code, 
        flake, 
        language_dependencies, 
        system_dependencies, 
        cmd_line_args, 
        compile_args,
        files,
        input,
        command,
        setup,
        language_version
    )
values
    ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
returning id;

-- name: GetExecRequest :one
select * from exec_request where id = $1;

-- name: GetExecRequestByHash :one
select * from exec_request where hash = $1;

-- name: ListExecRequests :many
select * from exec_request
where id >= $1
limit $2;

-- name: DeleteExecRequest :exec
delete from exec_request where id = $1;