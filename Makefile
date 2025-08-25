include .env
export
MIGRATIONS_PATH=./migrate/migrations
BINARY_NAME=main

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
migration:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up:
	@go run migrate/main.go up

.PHONY: migrate-down
migrate-down:
	@go run migrate/main.go down

.PHONY: migrate-step
migrate-step:
	@go run migrate/main.go step $(filter-out $@,$(MAKECMDGOALS))