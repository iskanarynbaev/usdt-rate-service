include .env
export

APP_NAME=app

build:
	go build -o $(APP_NAME) ./cmd/app

docker-build:
	docker build -t yourusername/usdtrate .

run:
	./$(APP_NAME)

lint:
	golangci-lint run

migrate-up:
	goose -dir migrations postgres "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up

migrate-down:
	goose -dir migrations postgres "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" down

docker-run:
	docker-compose up --build -d

run-all: migrate-up docker-run

GOOGLEAPIS_PATH=$(shell go list -m -f '{{.Dir}}' github.com/googleapis/googleapis)

gen-proto:
	protoc -I. -I$(GOOGLEAPIS_PATH) \
		--go_out=paths=source_relative:./internal/pb \
		--go-grpc_out=paths=source_relative:./internal/pb \
		--grpc-gateway_out=paths=source_relative:./internal/pb \
		--openapiv2_out=./internal/pb \
		proto/rate.proto

.PHONY: test

test:
	go test ./... -v -race
