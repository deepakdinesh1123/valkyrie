CREATE TABLE IF NOT EXISTS Worker (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(50),
    modified_at TIMESTAMP,
    modified_by VARCHAR(50)
);

CREATE TABLE IF NOT EXISTS JobQueue(
    id bigserial PRIMARY KEY,
    created_by VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    script TEXT,
    args jsonb,
    logs TEXT,
    flake TEXT,
    language VARCHAR(50),
    mem_peak INT,
    timeout INT,
    priority INT,
    lease_timeout FLOAT,
    queue VARCHAR(50) NOT NULL DEFAULT 'default',
    job_type VARCHAR(50) NOT NULL DEFAULT 'execution',
    worker_id INT REFERENCES Worker(id)
);

