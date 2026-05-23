APP_NAME=devmate-ai

.PHONY: test run-cli run-api build-cli build-api docker-build docker-up docker-down

test:
    go test ./...

run-cli:
    go run ./cmd/devmate

run-api:
    go run ./cmd/api

build-cli:
    go build -o bin/devmate ./cmd/devmate

build-api:
    go build -o bin/devmate-api ./cmd/api

docker-build:
    docker build -t devmate-api:local .

docker-up:
    docker compose up --build

docker-down:
    docker compose down
