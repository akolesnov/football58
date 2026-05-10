GO_CACHE ?= /tmp/football58-go-cache

.PHONY: test up down migrate-up migrate-down migrate-version

test:
	cd backend && GOCACHE=$(GO_CACHE) go test ./...

up:
	docker compose up -d --build

down:
	docker compose down

migrate-up:
	docker compose --profile tools run --rm migrate up

migrate-down:
	docker compose --profile tools run --rm migrate down 1

migrate-version:
	docker compose --profile tools run --rm migrate version
