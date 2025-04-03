# Enable .env file loading
set dotenv-load := true

import 'schemas/schemas.just'
import 'builds/builds.just'

# Common variables
export PG_URL := "postgresql://" + env_var("POSTGRES_USER") + ":" + env_var("POSTGRES_PASSWORD") + "@" + env_var("POSTGRES_HOST") + ":" + env_var("POSTGRES_PORT") + "/" + env_var("POSTGRES_DB") + "?sslmode=" + env_var("POSTGRES_SSL_MODE")

# Default values
export POSTGRES_USER := env_var_or_default("POSTGRES_USER", "thors")
export POSTGRES_PASSWORD := env_var_or_default("POSTGRES_USER", "thorkell")
export POSTGRES_DB := env_var_or_default("POSTGRES_DB", "valkyrie")
export POSTGRES_PORT := env_var_or_default("POSTGRES_PORT", "5432")
export POSTGRES_HOST := env_var_or_default("POSTGRES_HOST", "localhost")
export POSTGRES_SSL_MODE := env_var_or_default("POSTGRES_SSL_MODE", "disable")

# Run database migrations up
migrate:
    @echo "Running database migrations..."
    migrate -path internal/db/migrations -database "${PG_URL}" up

# Run database migrations down
migrate-down:
    @echo "Running database migrations down..."
    migrate -path internal/db/migrations -database "${PG_URL}" down

# Generate SQL code
gq:
    @echo "Generating SQL code..."
    cd internal/db && sqlc generate

# Start the PostgreSQL database
start-db:
    @echo "Starting PostgreSQL database..."
    podman run -d \
        --name postgres-container \
        -e POSTGRES_DB="${POSTGRES_DB}" \
        -e POSTGRES_USER="${POSTGRES_USER}" \
        -e POSTGRES_PASSWORD="${POSTGRES_PASSWORD}" \
        -p ${POSTGRES_PORT}:5432 \
        postgres

# Clear the state database
clear-stdb:
    @echo "Clearing state database..."
    rm -rf ~/.valkyrie/data

# Start observability services
start-observability: 
    @echo "Starting observability services..."
    docker-compose up -d valkyrie-otel-collector jaeger prometheus

# Build the valkyrie binary
valkyrie:
    @echo "Building valkyrie binary..."
    go build -o valkyrieb -tags docker cmd/valkyrie/main.go

# Start PostgreSQL with Docker and run migrations
docker-db: 
    @echo "Starting PostgreSQL and running migrations..."
    docker-compose up -d postgres
    sleep 5 # Wait for PostgreSQL to be ready
    just migrate

update-nix-hash nix_file package_name:
    hack/update_nix_hash.sh {{nix_file}} {{package_name}}

# Help command to list all available commands
help:
    @just --list

# Default recipe lists all available recipes
default:
    @just --list
