CREATE TABLE ExecutionQueue(
    id SERIAL PRIMARY KEY,
    worker_id VARCHAR(50) NOT NULL,
    created_by VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    started_at TIMESTAMP,
    executing BOOLEAN DEFAULT false,
    script TEXT,
    args JSON,
    logs TEXT,
    canceled BOOLEAN DEFAULT false,
    flake TEXT,
    language VARCHAR(50),
    mem_peak INT,
    timeout INT,
    priority INT,
    same_worker BOOLEAN DEFAULT false
);

