set dotenv-load := true

mod go 'just/go.just'
mod web 'just/web.just'

# List available recipes
default:
    @just --list

# Install web dependencies and git hooks
install:
    just web install
    lefthook install

# Format all code (Go + Web)
fmt:
    just go fmt
    just web fmt

# Lint all code (Go + Web)
lint:
    just go lint
    just web lint

# Run all test suites (Go + Web)
test:
    just go test
    just web test

# Build complete single binary with embedded frontend
build:
    just web build
    mkdir -p bin
    go build -ldflags="-s -w" -o bin/wwb ./cmd/server

# Generate OpenAPI / Swagger docs
swagger:
    swag init -g cmd/server/main.go -o internal/api/docs

# Clean build artifacts and local sqlite databases
clean:
    rm -rf bin/ web/dist/ *.db *.db-wal *.db-shm
