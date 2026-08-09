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

# Run Go development server
dev:
    go run ./cmd/server
