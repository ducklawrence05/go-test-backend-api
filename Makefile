include .env
export

BINARY_NAME=main
MIGRATIONS_PATH=./migrations
WIRE_PATH=./internal/infrastructure/wire

.PHONY: build
build:
	@go build -o ./bin/$(BINARY_NAME).exe cmd/app/main.go
 
.PHONY: run
run: build
	@./bin/$(BINARY_NAME).exe

.PHONY: test
test:
	@go test -v ./...

.PHONY: migrate-create
migrate-create:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up:
	@migrate -path $(MIGRATIONS_PATH) -database "$(POSTGRES_DB_URL)" up

.PHONY: migrate-down
migrate-down:
	@migrate -path $(MIGRATIONS_PATH) -database "$(POSTGRES_DB_URL)" down $(filter-out $@,$(MAKECMDGOALS))

.PHONY: wire
wire:
	@wire $(if $(filter-out $@,$(MAKECMDGOALS)),$(WIRE_PATH)/$(filter-out $@,$(MAKECMDGOALS)),$(WIRE_PATH)/...)
