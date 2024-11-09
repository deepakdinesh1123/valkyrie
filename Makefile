# Include other Makefiles
include .env oas/Makefile build/Makefile testing/Makefile

# PostgreSQL Connection URL
POSTGRES_URL = postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=${POSTGRES_SSL_MODE}

# Phony Targets
.PHONY: all migrate gq start-db clear-stdb start-observability odin docker-db add-pkgs run-pkgs store-pkgs dump

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
	@podman run -d \
		--name postgres-container \
		-e POSTGRES_DB=${POSTGRES_DB} \
		-e POSTGRES_USER=${POSTGRES_USER} \
		-e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} \
		-p 5432:5432 \
		postgres

# Clear the state database
clear-stdb:
	@echo "Clearing state database..."
	@rm -rf ~/.zango/data

# Start observability services
start-observability:
	@echo "Starting observability services..."
	@docker-compose up -d valkyrie-otel-collector jaeger prometheus

# Build the odin binary
odin:
	go build -o odinb cmd/odin/main.go

# Start PostgreSQL with Docker and run migrations
docker-db:
	@docker-compose up -d postgres 
	migrate -path internal/odin/db/migrations -database $(POSTGRES_URL) up

# Add packages from a dump file
add-pkgs:
	@if [ -z "$(filter-out $@,$(MAKECMDGOALS))" ]; then \
		echo "Error: Please provide the dump file name as an argument."; \
		exit 1; \
	fi; \
	dump_file=$(filter-out $@,$(MAKECMDGOALS)); \
	dump_path=./dumps/$$dump_file; \
	if [ ! -f "$$dump_path" ]; then \
		echo "Error: Dump file '$$dump_file' does not exist in the dumps folder."; \
		exit 1; \
	fi; \
	psql ${POSTGRES_URL} -c "DROP TABLE IF EXISTS packages CASCADE"; \
	echo "Applying $$dump_file to database..."; \
	psql ${POSTGRES_URL} -f $$dump_path; \
	psql ${POSTGRES_URL} -c "UPDATE packages SET tsv_search = to_tsvector('english', COALESCE(name, '') || ' ' || COALESCE(version, '') || ' ' || COALESCE(language, ''));"; \
	echo "Full-text search vectors generated successfully."; \
	psql ${POSTGRES_URL} -c "CREATE INDEX IF NOT EXISTS idx_packages_tsv ON packages USING GIN(tsv_search);"; \
	echo "GIN index for tsv_search created successfully."

# Store packages in the database
store-pkgs:
	@hack/store-pkgs.sh

# Dump package versions
dump:
	@if [ -z "$(filter-out $@,$(MAKECMDGOALS))" ]; then \
		echo "Error: Please specify a version (e.g., make dump 24.05)"; \
		exit 1; \
	fi
	@hack/packages.sh $(filter-out $@,$(MAKECMDGOALS))

# Catch-all target to suppress errors for non-existent targets
%:
	@:
