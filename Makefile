
VERSION ?= $(shell git describe --tags --always)
export LDFLAGS += -X github.com/enrichman/kubectl-epinio/internal/cmd.Version=$(VERSION)

build:
	go build -v -ldflags '$(LDFLAGS)' -o output/ ./...

lint:
	golangci-lint run

test:
	go test -v -race -covermode=atomic -coverprofile=coverage.out ./...
