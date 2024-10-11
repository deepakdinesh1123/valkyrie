#!/bin/bash

# Run PostgreSQL container
docker run -d \
  --name nixos-packages-db \
  -e POSTGRES_DB=nixos_packages \
  -e POSTGRES_USER=thors \
  -e POSTGRES_PASSWORD=nixpassword \
  -p 5433:5432 \
  -v $PWD/dumps:/dumps \
  postgres:16

# Wait for the database to be ready
sleep 10

# Create the table
docker exec -it nixos-packages-db psql -U thors -d nixos_packages -c "
CREATE TABLE IF NOT EXISTS packages (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    version VARCHAR(255) NOT NULL,
    type VARCHAR(255),
    language VARCHAR(255),
    arch VARCHAR(255),
    path TEXT NOT NULL
);"

echo "PostgreSQL container is set up and running on port 5433"

# Run the Nix build and nixdump command
# nix build
./result/bin/odin nixdump -c $1

# Destroy the container
docker stop nixos-packages-db
docker rm nixos-packages-db

rm store-paths.xz


