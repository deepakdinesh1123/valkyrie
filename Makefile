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
	@podman run -d \
		--name postgres-container \
		-e POSTGRES_DB=${POSTGRES_DB} \
		-e POSTGRES_USER=${POSTGRES_USER} \
		-e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} \
		-p 5432:5432 \
		postgres

.PHONY: migrate-db
migrate-db:
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
		-t odin:alpine \
		-f build/platforms/alpine/Dockerfile .

.PHONY: build-docker-image-ubuntu
build-docker-image-ubuntu:
	docker build \
		-t odin:ubuntu \
		--build-arg HOST_UID=$(shell id -u) \
		--build-arg HOST_GID=$(shell id -g) \
		--build-arg HOST_USER=$(shell whoami) \
		--build-arg HOST_GROUP=$(shell whoami) \
		-f build/platforms/ubuntu/Dockerfile .

.PHONY: build-podman-image-ubuntu
build-podman-image-ubuntu:
	podman build \
		-t odin:ubuntu \
		--build-arg HOST_UID=$(shell id -u) \
		--build-arg HOST_GID=$(shell id -g) \
		--build-arg HOST_USER=$(shell whoami) \
		--build-arg HOST_GROUP=$(shell whoami) \
		-f build/platforms/ubuntu/Containerfile .

.PHONY: build-podman-image
build-podman-image:
	podman build \
		-t odin:alpine \
		-f build/platforms/alpine/Containerfile .

.PHONY: odin
odin:
	go build -o odinb cmd/odin/main.go

.PHONY: docker-db
docker-db:
	@docker compose up postgres -d 
	migrate -path internal/odin/db/migrations -database $(POSTGRES_URL) up

.PHONY: add-pkgs run-pkgs
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
	echo "GIN index for tsv_search created successfully.";

store-pkgs:
	@packages=$$(psql ${POSTGRES_URL} -t -c "SELECT id, name, language FROM packages;"); \
	if [ -z "$$packages" ]; then \
		echo "No packages found in the database."; \
		exit 0; \
	fi; \
	while IFS="|" read -r id name language; do \
		name=$$(echo $$name | xargs); \
		language=$$(echo $$language | xargs); \
		if [ -z "$$language" ]; then \
			echo "Running nix-shell for $$name (type: system)..."; \
			nix-shell -p $$name --run "exit"; \
		else \
			echo "Running nix-shell for $$language.$$name (type: language)..."; \
			nix-shell -p $$language.$$name --run "exit"; \
		fi; \
	done <<< "$$packages"; \
	echo "All packages processed successfully.";

.PHONY: dump
dump:
	@if [ -z "$(filter-out $@,$(MAKECMDGOALS))" ]; then \
		echo "Error: Please specify a version (e.g., make dump 24.05)"; \
		exit 1; \
	fi
	./hack/packages.sh $(filter-out $@,$(MAKECMDGOALS))

%:
	@: