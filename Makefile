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
		-f build/platforms/ubuntu.dockerfile .

.PHONY: build-podman-image
build-podman-image:
	podman build \
		-t odin:latest \
		-f build/platforms/ubuntu.dockerfile .

.PHONY: odin
odin:
	go build -o odinb cmd/odin/main.go