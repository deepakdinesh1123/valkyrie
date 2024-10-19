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