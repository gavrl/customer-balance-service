build:
	go build -o cmd/customer-api cmd/main.go

run-dev:
	docker-compose -f ./docker-compose.dev.yml up -d; ./cmd/rundev.sh

run:
	docker-compose up --build