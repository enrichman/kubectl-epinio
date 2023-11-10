
VERSION ?= $(shell git describe --tags --always)
export LDFLAGS += -X github.com/enrichman/kubectl-epinio/internal/cli/cmd.Version=$(VERSION)
export GOCOVERDIR=$(shell pwd)/output/coverage

build:
	go build -v -ldflags '$(LDFLAGS)' -o output/ ./...

infra-setup:
	./scripts/setup_cluster.sh

infra-teardown:
	k3d cluster delete epinio

lint:
	golangci-lint run

build-test-bin:
	go build -v -ldflags '$(LDFLAGS)' -cover -covermode=atomic -coverpkg ./... -o output/

test: test-unit test-integration

test-unit: build
	go test $(shell go list ./... | grep -v /tests) -v -race -covermode=atomic -coverprofile=coverage-unit.out -coverpkg ./...

test-integration: build-test-bin
	mkdir -p ${GOCOVERDIR} && rm -rf ${GOCOVERDIR}/*
	go test -v ./tests
	go tool covdata percent -i=output/coverage
	go tool covdata textfmt -i=output/coverage -o coverage-integration.out
