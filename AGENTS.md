# Agents Rules for World Wide Bulb (WWB)

This file specifies constraints and standards for AI Agents when working on World Wide Bulb.

## CRITICAL INSTRUCTIONS

1. **Master Reference**: You must read, understand, and strictly follow all instructions in the project root's coding guidelines:
   - **File Path**: [GUIDELINES.md](GUIDELINES.md)

2. **The Ponytail Workflow (Default Execution Mode)**:
   - Climb the senior dev ladder: YAGNI -> Reuse -> Stdlib -> Platform -> Deps -> Simplicity -> Minimal Diff.
   - Deletion over addition. Boring over clever. Fewest files possible.
   - Mark intentional simplifications with `// ponytail:` comments naming the known ceiling and upgrade path.
   - Non-trivial logic must leave behind at least one runnable check.

3. **Commit Policy**:
   - AI assistants are strictly prohibited from generating, suggesting, or running git commits or pushes under any circumstance. Commits are human-only.

4. **Plan Execution Policy**:
   - **Never** proceed to code changes or implementation commands directly from a plan unless the human developer explicitly replies with direct authorization (e.g., "go ahead", "approved", "execute").
   - A plan is strictly for review. Keep your hands off the files until permission is explicitly given.

5. **Comments & Standards**:
   - Strictly NO inline comments in function bodies explaining syntax or obvious actions.
   - Ceiling of 300 lines per file. Always ask before performing major file splits or autonomous refactoring.
   - Strictly modern Svelte 5 runes (`$state`, `$derived`, `$props`) on the frontend; no legacy v4 syntax.
   - Strictly standard Go layout, implicit interfaces, early returns, and structured `log/slog` on the backend.
   - Web package manager defaults strictly to `pnpm` or `bun`.

Always consult the full rules inside [GUIDELINES.md](GUIDELINES.md) before writing code.
