# Include other Makefiles
include .env oas/Makefile build/Makefile testing/Makefile

# PostgreSQL Connection URL
POSTGRES_URL = postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=${POSTGRES_SSL_MODE}

# Phony Targets
.PHONY: all migrate gq start-db clear-stdb start-observability odin

# Default target
all: odin migrate start-db start-observability

# Run database migrations
migrate:
	@echo "Running database migrations..."
	@migrate -path internal/odin/db/migrations -database ${POSTGRES_URL} up

# Generate SQL code
gq: start-db
	@echo "Generating SQL code..."
	@cd internal/odin/db && sqlc generate

# Start the PostgreSQL database
start-db:
	@echo "Starting PostgreSQL database..."
	@docker compose up postgres -d
	@migrate -path internal/odin/db/migrations -database ${POSTGRES_URL} up

# Clear the state database
clear-stdb:
	@echo "Clearing state database..."
	@rm -rf ~/.zango/data

# Start observability services
start-observability:
	@echo "Starting observability services..."
	@docker compose up valkyrie-otel-collector jaeger prometheus -d

# Build the Odin application
odin:
	@echo "Building the Odin application..."
	@go build -o odinb cmd/odin/main.go
