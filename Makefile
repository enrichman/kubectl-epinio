
VERSION ?= $(shell git describe --tags --always)
export LDFLAGS += -X github.com/enrichman/kubectl-epinio/internal/cmd.Version=$(VERSION)

build:
	go build -v -ldflags '$(LDFLAGS)' -o output/ ./...

infra-setup:
	./scripts/setup_cluster.sh

infra-teardown:
	k3d cluster delete epinio

lint:
	golangci-lint run

build-test-bin:
	go build -v -ldflags '$(LDFLAGS)' -cover -o output/ ./...

test:
	go test -v -race -covermode=atomic -coverprofile=coverage.out ./...
