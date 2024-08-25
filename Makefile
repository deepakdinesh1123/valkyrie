include .env

POSTGRES_URL = postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=${POSTGRES_SSL_MODE}

TEST_POSTGRES_USER = test
TEST_POSTGRES_PASSWORD = test
TEST_POSTGRES_HOST = localhost
TEST_POSTGRES_PORT = 5400
TEST_POSTGRES_DB = test
TEST_POSTGRES_SSL_MODE = disable
DELETE_TEST_DB ?= true

.PHONY: migrate
migrate:
	@migrate -path internal/odin/db/migrations -database ${POSTGRES_URL} up

.PHONY: gq
gq: start-db
	# cd internal/odin/db && sqlc verify
	cd internal/odin/db && sqlc generate

.PHONY: start-server
start-server:
	@go run cmd/odin/main.go server start

.PHONY: start-worker
start-worker:
	@go run cmd/odin/main.go worker start

.PHONY: standalone
standalone:
	@go run cmd/odin/main.go standalone

.PHONY: start-db
start-db:
	@docker compose up postgres -d
	migrate -path internal/odin/db/migrations -database ${POSTGRES_URL} up

.PHONY: odin
odin:
	@go build -o odin cmd/odin/main.go

.PHONY: clear-stdb
clear-stdb:
	rm -rf ~/.zango/data

.PHONY: oapi-gen
oapi-gen:
	@go generate

.PHONY: start-observability
start-observability:
	@docker compose up valkyrie-otel-collector jaeger prometheus -d

.PHONY: start-test-db
start-test-db:
	@docker run \
		--rm \
		--name odin-test-db \
		-e POSTGRES_USER=${TEST_POSTGRES_USER} \
		-e POSTGRES_PASSWORD=${TEST_POSTGRES_PASSWORD} \
		-e POSTGRES_DB=${TEST_POSTGRES_DB} \
		-p ${TEST_POSTGRES_PORT}:5432 \
		-d postgres

.PHONY: stop-test-db
stop-test-db:
	@docker stop odin-test-db

.PHONY: test
test: start-test-db
	@go test ./internal/odin/server
	if [ "${DELETE_TEST_DB}" = "true" ]; then \
		$(MAKE) stop-test-db; \
	fi