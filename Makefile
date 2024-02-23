SHELL := /bin/bash
BINARY_NAME ?= favsapi
MOCKS_DESTINATION=internal/mocks

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## build: build the application
.PHONY: build
build:
	go build -o ${BINARY_NAME} ./cmd/main

## run: run the application
.PHONY: run
run: build
	./$(BINARY_NAME)

## clean: clean all artifacts
.PHONY: clean
clean:
	go clean
	rm ${BINARY_NAME}

## test: run all tests
.PHONY: test
test:
	go test -v -race ./...

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out

## dep: download dependencies
.PHONY: dep
dep:
	go mod download

## lint: run linter
.PHONY: lint
lint:
	golangci-lint run -v -c golangci-lint.yml

## mocks: generate mocks
.PHONY: mocks
# put the files with interfaces you'd like to mock in prerequisites
# wildcards are allowed
mocks: internal/pkg/auth/interfaces.go
	@echo "Generating mocks..."
	@rm -rf $(MOCKS_DESTINATION)
	@for file in $^; do mockgen -source=$$file -destination=$(MOCKS_DESTINATION)/$${file#*/}; done

## docs: generate swagger 2.0 documentation
.PHONY: docs
docs:
	@echo "Generating docs"
	@swag init -g main.go -d cmd/main/,internal
