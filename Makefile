
build:
	go build -v -o output/ ./...

lint:
	golangci-lint run
