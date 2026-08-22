set dotenv-load := true

db_path := env_var_or_default('DB_PATH', 'bulb.db')

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

# Check formatting for all code (Go + Web)
fmt-check:
    just go fmt-check
    just web fmt-check

# Lint all code (Go + Web)
lint:
    just go lint
    just web lint

# Audit frontend dependencies
audit:
    just web audit


# Run all test suites (Go + Web)
test:
    just go test
    just web test

# Build complete single binary with embedded frontend
build:
    just web build
    mkdir -p bin
    go build -ldflags="-s -w" -o bin/wwb ./cmd/wwb

# Build static Linux amd64 binary for containers/production
build-prod:
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 just build


# Generate OpenAPI / Swagger docs
swagger:
    swag init -g cmd/wwb/main.go -o internal/api/docs

# Generate Markdown documentation for just commands
docgen *args='':
    go run ./cmd/docgen {{args}}

# Verify documentation is in sync with justfiles
check-docs:
    go run ./cmd/docgen -check


# Clean build artifacts and local sqlite databases
clean:
    rm -rf {{db_path}} {{db_path}}-wal {{db_path}}-shm
    find web/build -mindepth 1 ! -name '.gitkeep' -delete


