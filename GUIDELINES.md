# World Wide Bulb (WWB) — Coding Guidelines

This document defines the coding standards, architectural decisions, and workflows for World Wide Bulb. Every AI assistant and developer working on this project **must** adhere strictly to these rules.

---

## 1. Universal Principles (All Stacks)

### 1.1 The Ponytail Workflow (Lazy Senior Dev Mode)
Before writing any code, stop at the first rung that holds:
1. **YAGNI:** Does this need to be built at all?
2. **Reuse:** Does it already exist in this codebase? Reuse existing helpers/patterns.
3. **Stdlib First:** Does the Go or web standard library already do this? Use it.
4. **Native Platform:** Does a native browser or platform feature cover it? Use it.
5. **Existing Deps:** Does an already-installed dependency solve it? Use it.
6. **Simplicity:** Can this be one clean line? Make it one line.
7. **Minimal Diff:** Only then write the minimum code that works.

#### Ponytail Execution Rules:
- **No unrequested abstractions:** No boilerplate, layers, or interfaces nobody asked for.
- **Deletion over addition:** Boring over clever. Fewest files possible.
- **Shortcuts & Ceilings:** Mark intentional simplifications with a `// ponytail:` comment naming the ceiling and the upgrade path (e.g., `// ponytail: global mutex on in-memory map; upgrade to sharded map if >10k concurrent IPs`).
- **Non-Negotiables:** Never be lazy about input validation at trust boundaries, security, error handling preventing data loss, or explicit user requests. Non-trivial logic must have at least one runnable check.

### 1.2 Comments & Documentation
- **No Redundant Comments:** Absolutely no inline comments echoing code syntax (e.g., `count++ // increment count`).
- **If a comment doesn't answer "How?", "Why?", or "WTF?", omit it.**
- **Allowed Comments:**
  - Package-level documentation.
  - Standard doc comments on exported symbols (Godoc, JSDoc/TSDoc).
  - High-level `// why:` context explaining non-obvious design constraints.
  - Actionable `// todo:` comments explaining future improvements or known rough edges.
  - `// ponytail:` comments documenting intentional MVP simplifications.
- **Style:** Short comments in lowercase. Capital letters only for multi-line contextual explanations.
- **Preservation:** Always preserve existing Swagger/OpenAPI annotations and language doc conventions.

### 1.3 Function & File Limits
- **Single Responsibility (SOLID):** Functions must do one clear action.
- **DRY:** Never duplicate logic; generalize reusable utilities.
- **File Length Ceiling:** Strict maximum of **300 lines per file**. If a file approaches or exceeds this:
  - Decouple, group, and split into sub-modules or logical extensions.
  - Notify and request approval from the developer before performing major autonomous refactoring.
- **Clean Code:** Zero dead code, unused imports, or magic numbers (use named constants).

### 1.4 Development & Verification Workflows
- **Package Managers:** Web stack defaults strictly to `pnpm` or `bun`. Never run `npm` or `yarn`.
- **Pre-commit Automation:** Use `lefthook` for pre-commit checks.
- **AI Tooling Constraints:**
  - AI assistants must **never** run `git commit` or `git push`. Commits belong solely to the human.
  - AI assistants must **never** execute implementation plans without explicit human approval.
- **Commit Format:** Follow `action(part): description` convention (e.g., `feat(bulb): add cooldown limiter`).

---

## 2. Go (Backend API & Engine)

### 2.1 Tooling & Verification
- **Formatting:** Format using standard Go tooling:
  ```bash
  just go fmt
  ```
- **Linting:** Code must pass `gofmt` and `golangci-lint` default rules (`errcheck`, `govet`, `staticcheck`, `revive`):
  ```bash
  just go lint
  ```

### 2.2 Documentation & Endpoints
- **Swagger / OpenAPI:** Decorate all HTTP endpoints with `swaggo/swag` annotations (`@Summary`, `@Tags`, `@Produce`, `@Success`, `@Failure`, `@Router`).
- **Godoc:** Every exported package, struct, interface, and function must have standard Godoc documentation.

### 2.3 Structural Rules & Architecture
- **Standard Go Layout:**
  - `/cmd/server` -> Entry point (`main.go`).
  - `/internal` -> Private logic (`bulb`, `store`, `ws`, `api`). Absolute ban on external package exposure.
  - `/web` -> Embedded static assets via `//go:embed`.
- **Implicit Interfaces:** Rely on Go's implicit interfaces. Do not declare interfaces before consumer packages require them.
- **Error Handling:** Return early on errors (`if err != nil { return err }`). Avoid deep nesting.
- **Structured Logging:** Use `log/slog` for operational diagnostics. (Info = normal, Warn = client 4xx, Error = server 5xx). Never log sensitive PII or raw IP addresses.
- **Testing:**
  - Group subtests using `t.Run`.
  - Name test files `*_test.go`.
  - Use `stretchr/testify/assert` consistently.

---

## 3. Svelte 5 (Frontend SPA)

### 3.1 Svelte 5 Runes & Paradigm
- **Reactivity Paradigm:** Always use modern Svelte 5 runes (`$state`, `$derived`, `$props`, `$effect`, `$bindable`).
- **Strict Legacy Ban:** Absolutely no legacy Svelte v4 reactive declarations (`$:`) or `export let` syntax.
- **Component Design:** Atomic, decoupled UI components. UI layout remains independent of direct data-fetching states.
- **Styling:** Exclusively Tailwind CSS utility classes. Avoid inline style tags or un-scoped CSS blocks.

### 3.2 Structure & Tooling
- **Directory Structure:**
  ```text
  web/
  ├── src/
  │   ├── lib/
  │   │   └── components/   # Atomic UI elements (Bulb, HistoryLog, StepChart)
  │   ├── App.svelte        # Root SPA container
  │   └── main.ts           # Entry point
  ├── dist/                 # Static build target for Go embed
  └── package.json
  ```
- **Formatting & Linting:** `prettier` for code formatting, `eslint` for linting, and `svelte-check` for template safety.
- **Testing:** `vitest` for client unit/store tests.

---

## 4. Verification Commands Reference

| Action | Go (Backend) | Svelte (Frontend) |
| :--- | :--- | :--- |
| **Format** | `just go fmt` | `just web fmt` |
| **Lint / Analyze** | `just go lint` | `just web lint` |
| **Type Check** | — | `just web check` |
| **Test** | `go test ./...` | `just web test` |
| **Build** | `go build ./cmd/server` | `just web build` |
