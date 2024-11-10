#!/bin/bash
set -e

# Function to clean up resources
cleanup() {
    echo "Cleaning up resources..."
    docker stop nixos-packages-db 2>/dev/null || true
    docker rm nixos-packages-db 2>/dev/null || true
    rm -f nixpkgs_data.csv nixpkgs_data.json
}

# Set trap for cleanup on script exit
trap cleanup EXIT

# Run PostgreSQL container
echo "Starting PostgreSQL container..."
docker run -d \
    --name nixos-packages-db \
    -e POSTGRES_DB=nixos_packages \
    -e POSTGRES_USER=thors \
    -e POSTGRES_PASSWORD=nixpassword \
    -p 5433:5432 \
    -v "$PWD/dumps:/dumps" \
    postgres:16

# Wait for the database to be ready
echo "Waiting for the database to be ready..."
for i in {1..30}; do
    if docker exec nixos-packages-db pg_isready -U thors -d nixos_packages; then
        echo "Database is ready."
        break
    fi
    if [ $i -eq 30 ]; then
        echo "Timeout waiting for database to be ready. Exiting."
        exit 1
    fi
    sleep 1
done

# Create the table 
echo "Creating table..."
docker exec -i nixos-packages-db psql -U thors -d nixos_packages <<-EOSQL
CREATE TABLE IF NOT EXISTS packages (
    package_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    version VARCHAR(255) NOT NULL,
    pkgType VARCHAR(255),
    language VARCHAR(255),
    store_path VARCHAR(255),
    tsv_search TSVECTOR
);
EOSQL

echo "PostgreSQL container is set up and running on port 5433"

# Run the Nix build and nixdump command
echo "Running nixdump command..."
if [ -z "$1" ]; then
    echo "Error: No argument provided for nixdump command."
    exit 1
fi

./odinb nixdump -c "$1"

echo "Script completed successfully."