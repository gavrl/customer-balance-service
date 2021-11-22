ifneq (,$(wildcard ./.env))
    include .env
    export
endif


dev-up:
	docker-compose -f ./docker-compose.dev.yml up -d

dev-down:
	docker-compose -f ./docker-compose.dev.yml down --remove-orphans

dev-run:
	go run ./cmd/app/main.go

dev-run-race:
	go run -race ./cmd/app/main.go

up:
	docker-compose up

build:
	go build -o balance-api cmd/app/main.go

up-build:
	docker-compose up -d --build

down:
	docker-compose down --remove-orphans

migrate-up:
	migrate -path ./internal/store/pg/migrations -database postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:55432/${POSTGRES_DB}?sslmode=${POSTGRES_SSL_MODE} up

migrate-down:
	migrate -path ./internal/store/pg/migrations -database postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:55432/${POSTGRES_DB}?sslmode=${POSTGRES_SSL_MODE} down

migration:
	migrate create -ext sql -dir ./internal/store/pg/migrations -seq init