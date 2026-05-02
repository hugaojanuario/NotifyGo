.PHONY: run-api run-gateway build test up down down-clear

run-api:
	go run ./cmd/api/main.go

run-gateway:
	go run ./cmd/gateway/main.go

build:
	go build -o bin/api ./cmd/api/main.go
	go build -o bin/gateway ./cmd/gateway/main.go

test:
	go test -race ./...

up:
	docker-compose up

down:
	docker-compose down

down-clear:
	docker-compose down -v

migrate:
	migrate -path db/migrations -database "postgres://admin:123456@localhost:5432/notifygo?sslmode=disable" up
