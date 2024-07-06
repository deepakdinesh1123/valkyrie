CREATE TABLE Execution(
    id SERIAL PRIMARY KEY,
    execution_id UUID,
    code TEXT,
    environment TEXT,
    requested_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    result TEXT,
    execution_status varchar(20),
    executed_at TIMESTAMP
);

