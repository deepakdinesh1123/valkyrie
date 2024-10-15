include .env oas/Makefile

POSTGRES_URL = postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=${POSTGRES_SSL_MODE}

.PHONY: migrate
migrate:
	@migrate -path internal/odin/db/migrations -database ${POSTGRES_URL} up

.PHONY: gq
gq: start-db
	# cd internal/odin/db && sqlc verify
	cd internal/odin/db && sqlc generate

.PHONY: start-db
start-db:
	@docker compose up postgres -d
	migrate -path internal/odin/db/migrations -database ${POSTGRES_URL} up

.PHONY: clear-stdb
clear-stdb:
	rm -rf ~/.zango/data

.PHONY: start-observability
start-observability:
	@docker compose up valkyrie-otel-collector jaeger prometheus -d

.PHONY: build-docker-image
build-docker-image:
	docker build \
		-t odin:latest \
		--build-arg HOST_UID=$(shell id -u) \
		--build-arg HOST_GID=$(shell id -g) \
		--build-arg HOST_USER=$(shell whoami) \
		--build-arg HOST_GROUP=$(shell whoami) \
		-f build/platforms/ubuntu.dockerfile .

.PHONY: build-podman-image
build-podman-image:
	podman build \
		-t odin:latest \
		--build-arg HOST_UID=$(shell id -u) \
		--build-arg HOST_GID=$(shell id -g) \
		--build-arg HOST_USER=$(shell whoami) \
		--build-arg HOST_GROUP=$(shell whoami) \
		-f build/platforms/ubuntu.dockerfile .

.PHONY: odin
odin:
	go build -o odinb cmd/odin/main.go

.PHONY: podman-db
podman-db:
	@podman compose -f docker-compose.yml up   postgres -d 
	migrate -path internal/odin/db/migrations -database $(POSTGRES_URL) up

.PHONY: add-pkgs
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
	psql ${POSTGRES_URL} -f $$dump_path ; \
	psql ${POSTGRES_URL} -c "UPDATE packages SET tsv_search = to_tsvector('english', COALESCE(name, '') || ' ' || COALESCE(version, '') || ' ' || COALESCE(language, ''));"; \
	echo "Full-text search vectors generated successfully." ; \
	psql ${POSTGRES_URL} -c "CREATE INDEX IF NOT EXISTS idx_packages_tsv ON packages USING GIN(tsv_search);";\
	echo "GIN index for tsv_search created successfully." ;\