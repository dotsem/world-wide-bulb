<img width="200" alt="bulb on" src="web/src/lib/assets/bulb_on.svg" align="right">

# World Wide Bulb

[![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?style=flat&logo=go&logoColor=white)](https://go.dev/)
[![SvelteKit Version](https://img.shields.io/badge/SvelteKit-5-FF3E00?style=flat&logo=svelte&logoColor=white)](https://kit.svelte.dev/)
[![License](https://img.shields.io/badge/license-MIT-green?style=flat)](LICENSE)

World Wide Bulb (WWB) is a real-time, globally synchronized digital light switch. It exposes a single, shared boolean state (ON/OFF) live across all connected visitors.

Everything compiles into a **single standalone binary**. This makes it easy to host it yourself for fun.

>[!NOTE]
> This project is made for fun and doesn't really have a real world use, unless... 

[![Backend CI](https://github.com/dotsem/world-wide-bulb/actions/workflows/backend-ci.yml/badge.svg)](https://github.com/dotsem/world-wide-bulb/actions/workflows/backend-ci.yml)
[![Frontend CI](https://github.com/dotsem/world-wide-bulb/actions/workflows/frontend-ci.yml/badge.svg)](https://github.com/dotsem/world-wide-bulb/actions/workflows/frontend-ci.yml)
[![Build Single Binary CI](https://github.com/dotsem/world-wide-bulb/actions/workflows/build-bin-ci.yml/badge.svg)](https://github.com/dotsem/world-wide-bulb/actions/workflows/build-bin-ci.yml)
[![Documentation CI](https://github.com/dotsem/world-wide-bulb/actions/workflows/docs-ci.yml/badge.svg)](https://github.com/dotsem/world-wide-bulb/actions/workflows/docs-ci.yml)
[![Codecov](https://codecov.io/gh/dotsem/world-wide-bulb/branch/main/graph/badge.svg)](https://codecov.io/gh/dotsem/world-wide-bulb)


## Technical Architecture

- **Backend:** Go (gorilla/websocket, gin-gonic, sqlc)
- **Frontend:** Svelte 5 w/ Tailwind CSS
- **Database:** SQLite
- **Distribution:** Single standalone executable (`bin/wwb`)

## Prerequisites & Tooling

To build or develop World Wide Bulb, the following tools are required:

- **Go** (1.26+)
- **[`golangci-lint`](https://golangci-lint.run/)**: Fast Go linters runner for code quality and pre-commit validation.
- **node.js** & **pnpm** (for frontend asset compilation and development)
- **[`just`](https://github.com/casey/just)**: It is recommended to install [`just`](https://github.com/casey/just) for the best development experience. `just` acts as the primary task runner, unifying commands for building, linting, formatting, and running tests across Go and Svelte components.
- **[`lefthook`](https://github.com/evilmartians/lefthook)**: Required for automated pre-commit checks (`just install` hooks it up).

## Getting Started

### 1. Repository Setup

Clone the repository and copy the environment configuration:

```bash
git clone https://github.com/dotsem/world-wide-bulb.git
cd world-wide-bulb
cp .env.template .env
```

### 2. Dependency Installation

Use `just` to install web dependencies and Git hooks (lefthook):

```bash
just install
```

### 3. Local Development

Run the Go backend and Svelte frontend development servers concurrently:

```bash
just go dev

just web dev
```

### 4. Single Binary Build & Execution

To compile the static frontend assets and package them into the single Go binary:

```bash
just build
```

The resulting binary will be output to `bin/wwb`. You can execute it directly:

```bash
./bin/wwb
```

## Task Runner Reference (`just`)

<!-- JUST_COMMANDS_START -->

### Root Commands

| Command | Description |
| :--- | :--- |
| `just default` | List available recipes |
| `just install` | Install web dependencies and git hooks |
| `just fmt` | Format all code (Go + Web) |
| `just fmt-check` | Check formatting for all code (Go + Web) |
| `just lint` | Lint all code (Go + Web) |
| `just audit` | Audit frontend dependencies |
| `just test` | Run all test suites (Go + Web) |
| `just build` | Build complete single binary with embedded frontend |
| `just swagger` | Generate OpenAPI / Swagger docs |
| `just docgen` | Generate Markdown documentation for just commands |
| `just check-docs` | Verify documentation is in sync with justfiles |
| `just clean` | Clean build artifacts and local sqlite databases |

### Backend (`just go <cmd>`)

| Command | Description |
| :--- | :--- |
| `just go fmt` | Format Go code |
| `just go fmt-check` | Check Go formatting |
| `just go lint` | Lint Go code |
| `just go test` | Run Go tests |
| `just go test-ci` | Run Go tests for CI with gotestsum and coverage |
| `just go coverage` | Run Go tests with coverage report (excluding generated sqlc files) |
| `just go dev` | Run Go development server |

### Frontend (`just web <cmd>`)

| Command | Description |
| :--- | :--- |
| `just web install` | Install web dependencies |
| `just web fmt` | Format Web code |
| `just web fmt-check` | Check Web formatting |
| `just web lint` | Lint Web code |
| `just web check` | Run Svelte type-checking |
| `just web test` | Run Web tests |
| `just web build` | Build Svelte static assets |
| `just web audit` | Audit frontend dependencies |
| `just web dev` | Run frontend development server |

<!-- JUST_COMMANDS_END -->

## Contribution Guidelines

1. **Linting & Formatting:** All code must cleanly pass `just lint` and `just fmt`.
2. **Testing:** Verify that `just test` passes before opening pull requests.
3. **Frontend Rules:** Frontend components must strictly use modern Svelte 5 runes (`$state`, `$derived`, `$props`).
4. **Code Standards:** No redundant inline syntax comments; keep function implementations minimal, readable, and early-returning.

Read the full guidelines in [GUIDELINES](GUIDELINES.md).

## License

This project is licensed under the MIT License.
