SHELL := /bin/bash -o pipefail

VERSION ?= $(shell ./scripts/gitversion.sh)

DB_NAME=app_db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_DSN="postgres://$(DB_USER):$(DB_PASSWORD)@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=disable"
MIGRATIONS_DIR="./assets/migrations"

CREATE_DB_IF_NOT_EXISTS="SELECT 'create database $(DB_NAME)' where not exists (select from pg_database where datname = '$(DB_NAME)')\gexec"
DROP_DB_IF_EXISTS="SELECT 'drop database $(DB_NAME)' where exists (select from pg_database where datname = '$(DB_NAME)')\gexec"

.PHONY: help
help:
	@echo "Usage: make <TARGET>"
	@echo ""
	@echo "Available targets are:"
	@echo ""
	@echo "    create-db                   Creates the application database if it does not already exist"
	@echo ""
	@echo "    drop-db                     Drops the application database if it exists"
	@echo ""
	@echo "    new-migration               Creates a new migration file. You need to pass the name of the migration to this target. Example: make new-migration name='init_schema'."
	@echo ""
	@echo "    migrate-up                  Applies all up migrations"
	@echo ""
	@echo "    migrate-one-up              Applies 1 up migrations"
	@echo ""
	@echo "    migrate-down                Applies all down migrations"
	@echo ""
	@echo "    migrate-one-down            Applies 1 down migrations"
	@echo ""
	@echo ""
	@echo "    vendor                      Tidies the dependency packages and updates the vendor folder"
	@echo ""
	@echo "    docs                        Generates the swagger documentation for the API"
	@echo ""
	@echo "    build                       Build binary"
	@echo ""
	@echo "    run                         Run main process"
	@echo ""
	@echo "    dev                         Run main process and setups the dev environment"
	@echo ""
	@echo "    test-all                    Run all tests and report coverage"

.PHONY: create-db
create-db:
	@echo "creating database if it does not already exists"
	@echo "database dsn: $(DB_DSN)"
	@echo $(CREATE_DB_IF_NOT_EXISTS) | psql -h localhost -p $(DB_PORT) -U $(DB_USER)

.PHONY: drop-db
drop-db:
	@echo "dropping database if it exists"
	@echo "database dsn: $(DB_DSN)"
	@echo $(DROP_DB_IF_EXISTS) | psql -h localhost -p $(DB_PORT) -U $(DB_USER)

.POHNY: setup-test-db
setup-test-db:
	@echo "creating test database if it does not already exists"
	@echo "database dsn: $(TEST_DB_DSN)"
	@echo $(CREATE_TEST_DB_IF_NOT_EXISTS) | psql -h localhost -p $(DB_PORT) -U $(DB_USER)
	@echo "running migrations"
	@migrate -path $(MIGRATIONS_DIR) -database $(TEST_DB_DSN) -verbose up

.PHONY: new-migration
new-migration:
	@migrate create -ext sql -seq -dir $(MIGRATIONS_DIR) -seq $$name

.PHONY: migrate-up
migrate-up:
	@migrate -path $(MIGRATIONS_DIR) -database $(DB_DSN) -verbose up

.PHONY: migrate-one-up
migrate-one-up:
	@migrate -path $(MIGRATIONS_DIR) -database $(DB_DSN) -verbose up

.PHONY: migrate-down
migrate-down:
	@migrate -path $(MIGRATIONS_DIR) -database $(DB_DSN) -verbose down

.PHONY: migrate-one-down
migrate-one-down:
	@migrate -path $(MIGRATIONS_DIR) -database $(DB_DSN) -verbose down 1

.PHONY: vendor
vendor:
	@go mod tidy
	@go mod vendor

.PHONY: docs
docs:
	@swag init --parseInternal --parseDepth 1
	@swag fmt ./...

.PHONY: build
build:
	@echo "building binaries for version: $(VERSION)"
	@go build -ldflags "-w -s -X github.com/amahdian/ai-assistant-be/version.AppVersion=$(VERSION) -X github.com/amahdian/ai-assistant-be/version.GitVersion=$(VERSION)" -mod vendor -o ./build/app-bin ./
	@echo "generated binary file: app-bin"

.PHONY: run
run:
	@go run main.go serve

.PHONY: dev
dev: gen
	@go run main.go serve

.PHONY: test-all
test-all:
	@go test -p 1 -mod=vendor -coverpkg=./... -coverprofile=.testCoverage.txt ./...
	@go tool cover -func=.testCoverage.txt