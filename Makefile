BIN_CMD_BOT := "./bin/telegram_bot_service"
BIN_CMD_STORAGE := "./bin/cloud_storage_service"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN_CMD_BOT) -ldflags "$(LDFLAGS)" ./cmd/telegram_bot_service
	go build -v -o $(BIN_CMD_STORAGE) -ldflags "$(LDFLAGS)" ./cmd/cloud_storage_service  

run: build
	$(BIN_CMD_BOT) -config ./configs/telegram_bot_config.toml

run-storage: build
	$(BIN_CMD_STORAGE) -config ./configs/cloud_storage_config.toml

dbuild:
	docker compose -f ./deployments/docker-compose.yaml build --no-cache

up:
	docker compose -f ./deployments/docker-compose.yaml up

down:
	docker compose -f ./deployments/docker-compose.yaml down --rmi all

version: build
	$(BIN_CMD_BOT) version
	$(BIN_CMD_STORAGE) version

test:
	go test -race ./...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint
