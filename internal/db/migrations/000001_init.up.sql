-- Create the base audit mixin table
CREATE TABLE fullauditmixin (
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by VARCHAR(255) DEFAULT '',
    modified_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    modified_by VARCHAR(255) DEFAULT ''
);

-- Create sequences
CREATE SEQUENCE languages_id_seq AS BIGINT;
CREATE SEQUENCE language_versions_id_seq AS BIGINT;
CREATE SEQUENCE exec_request_id_seq AS INT;
CREATE SEQUENCE workers_id_seq AS INT CYCLE;
CREATE SEQUENCE jobs_id_seq AS BIGINT;
CREATE SEQUENCE executions_id_seq AS BIGINT;
CREATE SEQUENCE sandboxes_id_seq AS BIGINT;

-- Languages table inheriting from fullauditmixin
CREATE TABLE languages (
    id BIGINT PRIMARY KEY DEFAULT nextval('languages_id_seq'),
    name TEXT NOT NULL UNIQUE,
    extension TEXT NOT NULL,
    monaco_language TEXT NOT NULL,
    template TEXT NOT NULL,
    is_disabled BOOLEAN NOT NULL DEFAULT FALSE,
    default_code TEXT NOT NULL
) INHERITS (fullauditmixin);

-- Language versions table inheriting from fullauditmixin
CREATE TABLE language_versions (
    id BIGINT PRIMARY KEY DEFAULT nextval('language_versions_id_seq'),
    language_id BIGINT NOT NULL REFERENCES languages (id) ON DELETE CASCADE,
    version TEXT NOT NULL,
    nix_package_name TEXT,
    template TEXT,
    default_version BOOLEAN NOT NULL DEFAULT FALSE,
    is_disabled BOOLEAN NOT NULL DEFAULT FALSE,
    UNIQUE (language_id, nix_package_name)
) INHERITS (fullauditmixin);

-- Unique index for default version per language
CREATE UNIQUE INDEX unique_default_version_per_language
ON language_versions (language_id)
WHERE default_version = TRUE;

-- Exec request table inheriting from fullauditmixin
CREATE TABLE exec_request (
    id INT PRIMARY KEY DEFAULT nextval('exec_request_id_seq'),
    hash TEXT NOT NULL,
    code TEXT,
    flake TEXT NOT NULL,
    language_dependencies TEXT[],
    system_dependencies TEXT[],
    cmd_line_args VARCHAR(1024),
    compile_args VARCHAR(1024),
    files JSONB,
    input TEXT,
    command TEXT,
    setup TEXT,
    system_setup TEXT,
    pkg_index TEXT,
    extension TEXT,
    language_version BIGINT NOT NULL REFERENCES language_versions(id) ON DELETE SET NULL,
    secrets BYTEA
) INHERITS (fullauditmixin);

-- Workers table inheriting from fullauditmixin
CREATE TABLE workers (
    id INT PRIMARY KEY DEFAULT nextval('workers_id_seq'),
    name TEXT NOT NULL UNIQUE,
    last_heartbeat TIMESTAMPTZ,
    current_state TEXT NOT NULL CHECK (current_state IN ('active', 'stale')) DEFAULT 'active'
) INHERITS (fullauditmixin);

-- Jobs table inheriting from fullauditmixin
CREATE TABLE jobs (
    job_id BIGINT PRIMARY KEY DEFAULT nextval('jobs_id_seq'),
    updated_at TIMESTAMPTZ,
    time_out INT,
    started_at TIMESTAMPTZ,
    arguments JSONB,
    current_state TEXT NOT NULL CHECK (current_state IN ('pending', 'scheduled', 'completed', 'failed', 'cancelled')) DEFAULT 'pending',
    retries INT DEFAULT 0,
    max_retries INT DEFAULT 5,
    worker_id INT REFERENCES workers ON DELETE SET NULL,
    job_type TEXT NOT NULL CHECK (job_type IN ('execution', 'sandbox'))
) INHERITS (fullauditmixin);

-- Index on arguments column
CREATE INDEX arguments_idx ON jobs USING GIN (arguments);

-- Executions table inheriting from fullauditmixin
CREATE TABLE executions (
    exec_id BIGINT PRIMARY KEY DEFAULT nextval('executions_id_seq'),
    job_id BIGINT REFERENCES jobs ON DELETE SET NULL,
    worker_id INT REFERENCES workers ON DELETE SET NULL,
    started_at TIMESTAMPTZ NOT NULL,
    finished_at TIMESTAMPTZ NOT NULL,
    exec_request_id INT REFERENCES exec_request ON DELETE SET NULL,
    exec_logs TEXT NOT NULL,
    nix_logs TEXT,
    out_files BYTEA,
    success BOOLEAN
) INHERITS (fullauditmixin);

-- Sandboxes table inheriting from fullauditmixin
CREATE TABLE sandboxes (
    sandbox_id BIGINT PRIMARY KEY DEFAULT nextval('sandboxes_id_seq'),
    worker_id INT REFERENCES workers ON DELETE SET NULL,
    started_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    sandbox_url TEXT,
    sandbox_agent_url TEXT,
    password BYTEA,
    config JSONB,
    details JSONB,
    current_state TEXT NOT NULL CHECK (current_state IN ('down', 'running', 'failed', 'stopped', 'creating', 'pending')) DEFAULT 'pending'
) INHERITS (fullauditmixin);

-- Indexes on sandbox JSONB columns
CREATE INDEX sndbx_cnfg_idx ON sandboxes USING GIN (config);
CREATE INDEX sndbx_details_idx ON sandboxes USING GIN (details);

-- System package filters table inheriting from fullauditmixin
CREATE TABLE system_package_filters (
    id SERIAL PRIMARY KEY,
    filter_type VARCHAR(10) NOT NULL CHECK (filter_type IN ('include', 'exclude')),
    package_string VARCHAR(255) NOT NULL
) INHERITS (fullauditmixin);

-- Index on system package filters
CREATE INDEX idx_systempackagefilters_type_string ON system_package_filters (filter_type, package_string);

-- Users table inheriting from fullauditmixin
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE,
    email VARCHAR(255) UNIQUE
) INHERITS (fullauditmixin);

-- Insert initial data for languages
INSERT INTO languages (name, extension, monaco_language, template, default_code, created_by, modified_by) VALUES (
    'generic',
    'none',
    'plaintext',
    ' ',
    ' ',
    'system',
    'system'
);

-- Insert initial data for language_versions
INSERT INTO language_versions (language_id, version, nix_package_name, template, default_version, created_by, modified_by) VALUES (
    (SELECT id FROM languages WHERE name = 'generic'),
    '0.0.1',
    NULL,
    ' ',
    TRUE,
    'system',
    'system'
);

-- Function to automatically update modified_at timestamp
CREATE OR REPLACE FUNCTION update_modified_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.modified_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for all tables to auto-update modified_at
CREATE TRIGGER update_languages_modified_at BEFORE UPDATE ON languages
    FOR EACH ROW EXECUTE FUNCTION update_modified_at_column();

CREATE TRIGGER update_language_versions_modified_at BEFORE UPDATE ON language_versions
    FOR EACH ROW EXECUTE FUNCTION update_modified_at_column();

CREATE TRIGGER update_exec_request_modified_at BEFORE UPDATE ON exec_request
    FOR EACH ROW EXECUTE FUNCTION update_modified_at_column();

CREATE TRIGGER update_workers_modified_at BEFORE UPDATE ON workers
    FOR EACH ROW EXECUTE FUNCTION update_modified_at_column();

CREATE TRIGGER update_jobs_modified_at BEFORE UPDATE ON jobs
    FOR EACH ROW EXECUTE FUNCTION update_modified_at_column();

CREATE TRIGGER update_executions_modified_at BEFORE UPDATE ON executions
    FOR EACH ROW EXECUTE FUNCTION update_modified_at_column();

CREATE TRIGGER update_sandboxes_modified_at BEFORE UPDATE ON sandboxes
    FOR EACH ROW EXECUTE FUNCTION update_modified_at_column();

CREATE TRIGGER update_system_package_filters_modified_at BEFORE UPDATE ON system_package_filters
    FOR EACH ROW EXECUTE FUNCTION update_modified_at_column();

CREATE TRIGGER update_users_modified_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_modified_at_column();
