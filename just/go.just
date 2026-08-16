set working-directory := '..'

# Format Go code
fmt:
    gofmt -l -w -s .

# Lint Go code
lint:
    golangci-lint run --config=.golangci.yml

# Run Go tests
test:
    go test ./...

# Run Go tests for CI with gotestsum and coverage
test-ci:
    #!/usr/bin/env bash
    gotestsum --format pkgname-and-test-fails --junitfile unit-tests.xml -- -race ./... -coverprofile=coverage.out -covermode=atomic
    grep -v -E "internal/store/(db|models|.*\.sql)\.go" coverage.out > coverage_clean.out
    mv coverage_clean.out coverage.out

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


