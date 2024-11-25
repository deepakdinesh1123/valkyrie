create table sandboxes (
    sandbox_id bigint primary key default nextval('sandboxes_id_seq'),
    worker_id int references workers on delete set null,
    started_at timestamptz,
    updated_at timestamptz,
    created_at timestamptz not null default now(),
    git_url text,
    sandbox_url text,
    password bytea,
    current_state TEXT NOT NULL CHECK (current_state IN ('pending', 'running', 'failed', 'stopped', 'creating'))
);

-- name: InsertSandbox :one
insert into sandboxes (git_url)
values ($1)
returning *;

-- name: GetSandbox :one
select *
from sandboxes
where  sandbox_id = $1;

-- name: UpdateSandboxState :exec
update  sandboxes set
current_state = $2
where sandbox_id =  $1;

-- name: UpdateSandboxStartTime :exec
update sandboxes set
started_at = $2
where sandbox_id = $1;

-- name: MarkSandboxRunning :exec
update sandboxes set
started_at = now(),
sandbox_url = $2,
password = $3,
current_state = 'running',
updated_at = now()
where sandbox_id = $1;

-- name: UpdateSandboxPassword :exec
update sandboxes set
password = $2
where sandbox_id = $1;

-- name: FetchSandboxJob :one
update sandboxes set current_state = 'creating', started_at = now(), worker_id = $1, updated_at = now()
where sandbox_id = (
    select sandbox_id from sandboxes
    where 
        current_state = 'pending'
    order by
        sandbox_id asc
    for update skip locked
    limit 1
    )
returning *;


