create table exec_request (
    id int primary key default nextval('exec_request_id_seq'),
    hash text not null,
    code text,
    flake text not null,
    language_dependencies text[],
    system_dependencies text[],
    cmd_line_args varchar(1024),
    compile_args varchar(1024),
    files bytea,
    input text,
    command text,
    setup text,
    programming_language text not NULL default "bash"
);

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
        programming_language
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
limit $1 offset $2;

-- name: DeleteExecRequest :exec
delete from exec_request where id = $1;