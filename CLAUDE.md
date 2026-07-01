# CLAUDE.md

The primary instruction file for this repository is **[AGENTS.md](./AGENTS.md)** — read that first. It applies equally to Claude Code, Codex, Cursor, Copilot, Aider, and human contributors, and covers commands, architecture, code style, commit conventions, and the review checklist.

## Claude Code specifics

- **Verify before claiming done.** Run `make test` and `make format`, then check the diff with `git status` / `git diff` before reporting work as complete. Type-checking ≠ feature correctness.
- **Don't commit `.claude/`, scratch specs, or planning docs.** Keep them local. The repo's `.gitignore` does not currently exclude `.claude/` — be careful with `git add -A`.
- **Prefer `make` targets over raw `go` commands.** `make test` runs `go mod tidy` first and uses `-race`; ad-hoc `go test ./...` skips both.
- **Match existing patterns when extending the SDK.** When adding a service method, copy the shape of the nearest existing wrapper — the codebase is intentionally consistent (see *Service pattern* and *Method-pair convention* in AGENTS.md).
- **For breaking changes**, use a `feat!:` / `fix!:` commit with a migration note in the body. The CI Conventional Commits check is strict and will block merges.
- **v1beta is the default; `v1alpha/` is frozen.** New work goes in the root package. Touch `v1alpha/` only for bugfixes.
- **Never edit `go-codegen` generated types** — those changes belong upstream in the spec / codegen repo.

Everything else lives in [AGENTS.md](./AGENTS.md).
