CREATE TABLE ExecutionResult (
    id SERIAL PRIMARY KEY,
    result TEXT,
    cmdLineArgs TEXT,
    environment TEXT,
    flake TEXT,
    dependency_file TEXT,
    executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    sandbox UUID NOT NULL UNIQUE
);

