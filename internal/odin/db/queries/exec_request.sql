create table exec_request (
    id int primary key default nextval('exec_request_id_seq'),
    hash text not null,
    code text not null,
    path text not null,
    flake text not null,
    args varchar(1024),
    programming_language text
);

-- name: InsertExecRequest :one
insert into exec_request
    (hash, code, path, flake, args, programming_language, nix_script)
values
    ($1, $2, $3, $4, $5, $6, $7)
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