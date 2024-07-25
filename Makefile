include .env

POSTGRES_URL = postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=${POSTGRES_SSL_MODE}

.PHONY: migrate
migrate:
	@migrate -path internal/odin/db/migrations -database ${POSTGRES_URL} up

.PHONY: gq
gq:
	@cd internal/odin/db && sqlc generate

.PHONY: start-server
start-server:
	@go run cmd/odin/main.go server start

.PHONY: start-worker
start-worker:
	@go run cmd/odin/main.go worker start

.PHONY: start-db
start-db:
	@docker compose up postgres -d
	migrate -path internal/odin/db/migrations -database ${POSTGRES_URL} up

.PHONY: odin
odin:
	@go build -o odin cmd/odin/main.go

