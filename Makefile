include .env.tests
export $(shell sed 's/=.*//' .env.tests)

CURRENT_DIR = $(shell pwd)
LOCAL_BIN=$(CURRENT_DIR)/bin

run:
	go run ./cmd/server/main.go

swag:
	swag init -g cmd/server/main.go

test:
	APP_LEVEL=test go test -v ./... -coverprofile=cover.out

cover:
	go tool cover -html=cover.out

mocks:
	mockgen -source=./internal/service/todo/interfaces.go -destination=./pkg/mocks/service/todo/mock_todo.go

lint:
	golangci-lint run ./... --timeout 60s

db\:up: bin-deps
	$(LOCAL_BIN)/goose -dir=migrations postgres "user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) sslmode=disable" up

db\:down: bin-deps
	$(LOCAL_BIN)/goose -dir=migrations postgres "user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) sslmode=disable" down

db\:create: bin-deps
	$(LOCAL_BIN)/goose -dir=migrations create $(NAME) sql

bin-deps:
	@mkdir -p bin
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.5.3
