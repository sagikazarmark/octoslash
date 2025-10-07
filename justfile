default:
    just --list

test:
    go test -shuffle on -race -v ./...

lint:
    golangci-lint run

fmt:
    golangci-lint fmt
