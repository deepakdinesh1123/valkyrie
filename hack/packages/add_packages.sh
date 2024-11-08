#!/bin/bash

POSTGRES_URL=postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=${POSTGRES_SSL_MODE}

# Check if a dump file name was provided as an argument
if [ -z "$1" ]; then
    echo "Error: Please provide the dump file name as an argument."
    exit 1
fi

dump_file="$1"

# Check if the dump file exists
if [ ! -f "$dump_file" ]; then  
    echo "Error: Dump file '$dump_file' does not exist in the dumps folder."
    exit 1
fi

# Drop the packages table if it exists
psql "${POSTGRES_URL}" -c "DROP TABLE IF EXISTS packages CASCADE"

# Apply the dump file to the database
echo "Applying $dump_file to database..."
psql "${POSTGRES_URL}" -f "$dump_file"

# Update the full-text search vectors
psql "${POSTGRES_URL}" -c "UPDATE packages SET tsv_search = to_tsvector('english', COALESCE(name, '') || ' ' || COALESCE(version, '') || ' ' || COALESCE(language, ''));"
echo "Full-text search vectors generated successfully."

# Create GIN index for tsv_search
psql "${POSTGRES_URL}" -c "CREATE INDEX IF NOT EXISTS idx_packages_tsv ON packages USING GIN(tsv_search);"
echo "GIN index for tsv_search created successfully."
