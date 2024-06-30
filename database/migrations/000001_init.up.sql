CREATE TABLE ExecutionResult (
    id SERIAL PRIMARY KEY,
    execution_id UUID,
    result TEXT,
    code TEXT,
    devenv TEXT,
    execution_status varchar(20),
    executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

