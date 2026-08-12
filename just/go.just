set working-directory := '..'

# Format Go code
fmt:
    gofmt -s -w .

# Lint Go code
lint:
    golangci-lint run ./...

# Run Go tests
test:
    go test ./...

# Run Go tests with coverage report (excluding generated sqlc files)
coverage:
    #!/usr/bin/env bash
    go test -coverprofile=coverage.out ./...
    grep -v -E "internal/store/(db|models|.*\.sql)\.go" coverage.out > coverage_clean.out
    mv coverage_clean.out coverage.out
    go tool cover -func=coverage.out

# Run Go development server
dev:
    go run ./cmd/wwb

