# Include other Makefiles
include .env oas/Makefile build/Makefile

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

migrate-down:
	@migrate -path internal/odin/db/migrations -database ${POSTGRES_URL} down

# Generate SQL code
gq:
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
	go build -o odinb -tags all cmd/odin/main.go

# Start PostgreSQL with Docker and run migrations
docker-db:
	@docker-compose up -d postgres
	migrate -path internal/odin/db/migrations -database $(POSTGRES_URL) up

# Add packages from a dump file
add-pkgs:
	bash hack/packages/add_packages.sh $(filter-out $@,$(MAKECMDGOALS))

# Store packages in the database
store-pkgs:
	bash hack/packages/store_packages.sh

# Dump package versions
dump:
	@if [ -z "$(filter-out $@,$(MAKECMDGOALS))" ]; then \
		echo "Error: Please specify a version (e.g., make dump 24.05)"; \
		exit 1; \
	fi
	./hack/packages/dump_packages.sh $(filter-out $@,$(MAKECMDGOALS))
