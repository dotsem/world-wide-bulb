set working-directory := '..'

# Format Go code
fmt:
    gofmt -l -w -s .

# Check Go formatting
fmt-check:
    #!/usr/bin/env bash
    set -eo pipefail
    UNFORMATTED=$(gofmt -s -l .)
    if [ -n "$UNFORMATTED" ]; then
        echo "The following Go files are not formatted:"
        echo "$UNFORMATTED"
        exit 1
    fi

# Lint Go code
lint:
    golangci-lint run --config=.golangci.yml

# Run Go tests
test:
    go test ./...

# Run Go tests for CI with gotestsum and coverage
test-ci:
    #!/usr/bin/env bash
    set -uo pipefail
    gotestsum --format pkgname-and-test-fails --junitfile unit-tests.xml -- -race ./... -coverprofile=coverage.out -covermode=atomic
    TEST_EXIT=$?
    if [ -f coverage.out ]; then
        grep -v -E "internal/store/(db|models|.*\.sql)\.go" coverage.out > coverage_clean.out || true
        mv coverage_clean.out coverage.out
    fi
    exit $TEST_EXIT

# Run Go tests with coverage report (excluding generated sqlc files)
coverage:
    #!/usr/bin/env bash
    set -eo pipefail
    go test -coverprofile=coverage.out ./...
    grep -v -E "internal/store/(db|models|.*\.sql)\.go" coverage.out > coverage_clean.out
    mv coverage_clean.out coverage.out
    go tool cover -func=coverage.out

# Run Go development server
dev:
    go run ./cmd/wwb


