-- name: InsertSandbox :one
insert into sandboxes (config, details)
values ($1, $2)
returning *;

-- name: GetSandbox :one
select *
from sandboxes
where  sandbox_id = $1;

-- name: UpdateSandboxState :exec
update  sandboxes set
current_state = $2,
details = $3
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
sandbox_agent_url = $4,
current_state = 'running',
updated_at = now()
where sandbox_id = $1;

-- name: UpdateSandboxPassword :exec
update sandboxes set
password = $2
where sandbox_id = $1;

-- name: ClearSandboxes :exec
delete from sandboxes;
