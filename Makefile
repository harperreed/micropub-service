# Makefile for Go project

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=server

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/server

test:
	$(GOTEST) -v ./... -coverprofile=coverage.out -covermode=atomic && go tool cover -html=coverage.out -o coverage.html

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/server
	./$(BINARY_NAME)

deps:
	$(GOGET) ./...
	$(GOMOD) tidy

.PHONY: all build test clean run deps
