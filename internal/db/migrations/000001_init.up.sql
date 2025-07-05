create sequence languages_id_seq as bigint;

CREATE TABLE languages (
    id bigint PRIMARY KEY DEFAULT nextval('languages_id_seq'),
    name TEXT NOT NULL UNIQUE,
    extension TEXT NOT NULL,
    monaco_language TEXT NOT NULL,
    template TEXT NOT NULL,
    is_disabled BOOLEAN NOT NULL DEFAULT false,
    default_code TEXT NOT NULL
);

create sequence language_versions_id_seq as bigint;

INSERT INTO languages (name, extension, monaco_language, template, default_code) VALUES (
    'generic',
    'none',
    'plaintext',
    ' ',
    ' '
);

CREATE TABLE language_versions (
    id bigint PRIMARY KEY DEFAULT nextval('language_versions_id_seq'),
    language_id BIGINT NOT NULL REFERENCES languages (id) ON DELETE CASCADE,
    version TEXT NOT NULL,
    nix_package_name TEXT,
    template TEXT,
    default_version BOOLEAN NOT NULL DEFAULT false,
    is_disabled BOOLEAN NOT NULL DEFAULT false,
    UNIQUE (language_id, nix_package_name)
);

CREATE UNIQUE INDEX unique_default_version_per_language
ON language_versions (language_id)
WHERE default_version = true;

INSERT INTO language_versions (language_id, version, nix_package_name, template, default_version) VALUES (
    (SELECT id FROM languages WHERE name = 'generic'),
    '0.0.1',
    NULL,
    ' ',
    TRUE
);

create sequence exec_request_id_seq as int;

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
    system_setup text,
    pkg_index text,
    extension text,
    language_version BIGINT NOT NULL REFERENCES language_versions(id) ON DELETE SET NULL,
    secrets bytea
);

create table workers (
    id int primary key,
    name text not null unique,
    created_at timestamptz not null default now(),
    last_heartbeat timestamptz,
    current_state TEXT NOT NULL CHECK (current_state IN ('active', 'stale')) DEFAULT 'active'
);

create sequence workers_id_seq as int cycle owned by workers.id;
alter table workers alter column id set default nextval('workers_id_seq');

create sequence jobs_id_seq as bigint;

create table jobs (
    job_id bigint primary key default nextval('jobs_id_seq'),
    created_at timestamptz not null  default now(),
    updated_at timestamptz,
    time_out int,
    started_at timestamptz,
    arguments jsonb,
    current_state TEXT NOT NULL CHECK (current_state IN ('pending', 'scheduled', 'completed', 'failed', 'cancelled')) DEFAULT 'pending',
    retries int default 0,
    max_retries int default 5,
    worker_id int references workers on delete set null,
    job_type TEXT NOT NULL CHECK (job_type IN ('execution', 'sandbox'))
);

CREATE INDEX arguments_idx ON jobs USING gin (arguments);

create sequence executions_id_seq as bigint;

create table executions (
    exec_id bigint primary key default nextval('executions_id_seq'),
    job_id bigint references jobs on delete set null,
    worker_id int references workers on delete set null,
    started_at timestamptz not null,
    finished_at timestamptz not null,
    created_at timestamptz not null default now(),
    exec_request_id int references exec_request on delete set null,
    exec_logs text not null,
    nix_logs text,
    success boolean
);

create sequence sandboxes_id_seq as bigint;

create table sandboxes (
    sandbox_id bigint primary key default nextval('sandboxes_id_seq'),
    worker_id int references workers on delete set null,
    started_at timestamptz,
    created_at timestamptz not null default now(),
    updated_at timestamptz,
    sandbox_url text,
    sandbox_agent_url text,
    password bytea,
    config jsonb,
    details jsonb,
    current_state TEXT NOT NULL CHECK (current_state IN ('down', 'running', 'failed', 'stopped', 'creating', 'pending')) default 'pending'
);

create index sndbx_cnfg_idx ON sandboxes USING gin (config);
create index sndbx_details_idx ON sandboxes USING gin (details);
