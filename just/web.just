set working-directory := '..'

# Install web dependencies
install:
    cd web && pnpm install

# Format Web code
fmt:
    cd web && pnpm exec prettier --write .

# Check Web formatting
fmt-check:
    cd web && pnpm exec prettier --check .

# Lint Web code
lint:
    cd web && pnpm exec eslint . --fix

# Run Svelte type-checking
check:
    cd web && pnpm exec svelte-check

# Run Web tests
test:
    cd web && pnpm test

# Build Svelte static assets
build:
    cd web && pnpm build

# Audit frontend dependencies
audit:
    cd web && pnpm audit --audit-level high

# Run frontend development server
dev:
    cd web && pnpm dev

