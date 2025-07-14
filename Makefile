build:
	go build -o app ./cmd/app

run:
	go run ./cmd/app

test:
	go test ./... -v

proto:
	protoc --go_out=internal/grpc --go-grpc_out=internal/grpc proto/*.proto

docker-build:
	docker build -t usdt-rate-service .

lint:
	golangci-lint run

env:
	cat .env
