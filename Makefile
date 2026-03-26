GO ?= go
PACKAGES ?= ./...
BINARY_NAME ?= tflint-ruleset-trailing-comma
COVERAGE_FILE ?= coverage.out

.PHONY: default build test coverage install

default: build

build:
	mkdir -p bin
	$(GO) build -o bin/$(BINARY_NAME) .

test:
	$(GO) test $(PACKAGES)

coverage:
	$(GO) test -covermode=atomic -coverprofile=$(COVERAGE_FILE) $(PACKAGES)
	$(GO) tool cover -func=$(COVERAGE_FILE)

install: build
	mkdir -p ~/.tflint.d/plugins
	install --mode +x ./bin/$(BINARY_NAME) ~/.tflint.d/plugins/$(BINARY_NAME)
