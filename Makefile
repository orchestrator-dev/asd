.DEFAULT_GOAL := build

VERSION ?= dev

build:
	go build -ldflags="-s -w -X main.version=$(VERSION)" -o dist/asd .

test:
	go test -race -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

lint:
	golangci-lint run ./...

vet:
	go vet ./...

release:
	goreleaser release --clean

snapshot:
	goreleaser release --snapshot --clean

install:
	cp dist/asd /usr/local/bin/asd

clean:
	rm -rf dist/ coverage.out
