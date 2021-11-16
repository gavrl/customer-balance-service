ifneq (,$(wildcard ./.env))
    include .env
    export
endif

build:
	go build -o balance-api cmd/app/main.go


dev-up:
	docker-compose -f ./docker-compose.dev.yml up -d

dev-run:
	go run ./cmd/app/main.go

dev-run-race:
	go run -race ./cmd/app/main.go

up-build:
	docker-compose up --build

up:
	docker-compose up

migrate-up:
	migrate -path ./internal/store/pg/migrations -database postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=${POSTGRES_SSL_MODE} up

migrate-down:
	migrate -path ./internal/store/pg/migrations -database postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=${POSTGRES_SSL_MODE} down

migration:
	migrate create -ext sql -dir ./internal/store/pg/migrations -seq init