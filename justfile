default:
    just --list

build:
    go build ./cmd/octoslash

test:
    go test -shuffle on -race -v ./...

lint:
    golangci-lint run

fmt:
    golangci-lint fmt
