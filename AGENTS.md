# AGENTS.md

Guidance for any coding agent or contributor working on this repository. This is the canonical instruction file — Claude Code, Codex, Cursor, Copilot, Aider, and human contributors should all start here. `CLAUDE.md` exists for Claude Code-specific overrides only and points back to this file.

## What this is

The UTxO RPC Go SDK — a Go client for [UTxO RPC](https://utxorpc.org) servers. Implements the UTxO RPC spec (Query, Submit, Sync, Watch services) over HTTP/2 using Connect RPC, with optional Cardano-specific helpers.

- Module: `github.com/utxorpc/go-sdk`
- Go: 1.25+ (see `go.mod`)
- License: Apache-2.0

## Repository layout

```
/                  Generic SDK package (`sdk`) — v1beta clients (default)
  main.go          UtxorpcClient, options, header injection
  query.go         Query service wrappers
  submit.go        Submit service wrappers
  sync.go          Sync service wrappers
  watch.go         Watch service wrappers
  doc.go           Package GoDoc
cardano/           Cardano-specific high-level client
v1alpha/           Legacy v1alpha clients (frozen, kept for back-compat)
examples/          Runnable examples (one binary per subdir)
.github/workflows/ CI: go-test, golangci-lint, conventional-commits, publish
```

## Common commands

```bash
make test       # mod-tidy + go test -v -race ./...
make format     # go fmt ./... + gofmt -s -w
make build      # builds every example binary
make mod-tidy   # go mod tidy
make golines    # reformat to 80-col lines (chain-split-dots)
```

There is no `make lint` target; CI runs `golangci-lint` directly. Run it locally with `golangci-lint run` if you want parity.

Examples need:

```bash
export UTXORPC_URL="https://preview.utxorpc-v0.demeter.run"
export DMTR_API_KEY="your-api-key"
go run examples/query/main.go
```

## Architecture

### Two layers

**Generic SDK** (root package `sdk`): blockchain-agnostic. `UtxorpcClient` exposes four service fields — `Query`, `Submit`, `Sync`, `Watch` — created via `NewClient(opts...)` using functional options (`WithBaseUrl`, `WithHeaders`, `WithHttpClient`, `WithDialTimeout`, `WithRequestTimeout`, etc.).

**Cardano layer** (`cardano` package): wraps `UtxorpcClient` with Cardano-specific helpers — hex/base64 decoding of tx hashes and addresses, methods like `GetUtxoByRef`, `SubmitTransaction`, `GetTip`. Built on top of `gouroboros` for address handling.

### Versioning

- **v1beta** is the default. New code uses `github.com/utxorpc/go-codegen/utxorpc/v1beta/...`.
- **v1alpha** lives under `./v1alpha/` and is frozen. Add bug fixes only — no new features. The default switch happened in commit `dab333c` (`feat!: default to v1beta and move v1alpha to sdk/v1alpha`).
- Don't introduce a third parallel layout. If a new spec version arrives, mirror the existing pattern.

### Service pattern

Every service file follows the same shape:

1. Type alias to the generated Connect client, e.g. `type QueryServiceClient = queryconnect.QueryServiceClient`.
2. A constructor method on `UtxorpcClient` that builds the Connect client with the configured transport and base URL.
3. Wrapper methods that call `AddHeadersToRequest()` to inject configured headers before dispatching.

Wrappers exist so users don't have to remember to attach headers on every request.

### Method-pair convention

Each user-facing wrapper comes in two forms:

- `MethodName(req)` — uses `context.Background()`.
- `MethodNameWithContext(ctx, req)` — caller supplies the context.

Always add both when introducing a new method. Don't ship just one half.

### Streaming methods

These return `*connect.ServerStreamForClient`:

- `FollowTip` — chain-tip changes (Apply / Undo / Reset).
- `WaitForTx` — transaction confirmation.
- `WatchMempool` — mempool changes.
- `WatchTx` — transaction watching.

Streaming wrappers also follow the method-pair convention.

### Key dependencies

- `connectrpc.com/connect` — RPC framework (gRPC-compatible, HTTP/2).
- `github.com/utxorpc/go-codegen` — generated protobuf types for the spec.
- `github.com/blinklabs-io/gouroboros` — Cardano primitives (address handling).
- `golang.org/x/net/http2` — explicit HTTP/2 transport.

## Code style

- **Formatting**: `gofumpt`, `gofmt -s`, `goimports`, `gci`. All declared in `.golangci.yml` and run by CI. Run `make format` (and optionally `make golines`) before committing.
- **Line length**: 80 columns when running `make golines`. Don't fight the tool — let it split chains on dots.
- **Linting**: CI runs `golangci-lint` with the set in `.golangci.yml` (`errorlint`, `gosec`, `perfsprint`, `prealloc`, `protogetter`, `bodyclose`, `contextcheck`, `errcheckjson`, `nilerr`, `unparam`, etc.). Match existing style; don't add `//nolint:` without a real reason.
- **Errors**: use `errors.Is` / `errors.As` (errorlint is on); never compare errors with `==`.
- **Performance**: pre-allocate slices when length is known (`prealloc`); prefer `strconv` over `fmt.Sprintf` for trivial conversions (`perfsprint`).
- **Protobuf accessors**: use generated `GetX()` getters rather than direct field access (`protogetter`).
- **Generated code**: never edit `go-codegen` types in vendored or local copies; spec changes happen upstream.

## Testing

- Tests run with `-race`. Don't introduce code paths with data races even if existing tests don't catch them.
- `make test` runs `go mod tidy` first — tidy diffs in PRs are a sign someone skipped this.
- Network/integration tests against a live UTxO RPC server are not part of `make test`. If you add one, gate it behind an env var or `testing.Short()`.
- Examples under `examples/` are not run by `make test`; `make build` only verifies they compile.

## Commits & PRs

- **Conventional Commits are enforced by CI** (`.github/workflows/conventional-commits.yml`). Use `feat:`, `fix:`, `chore:`, `docs:`, `refactor:`, `test:`, `ci:`, etc. PRs that fail this check cannot merge.
- Breaking changes use `feat!:` / `fix!:` and explain the migration in the commit body and PR description.
- Dependency bumps come through Dependabot — match its commit style (`chore(deps): bump …`) for any manual bumps.
- Run `make test` and `make format` before pushing.
- Keep commits focused; prefer multiple small commits over one large one when changes are logically separable.

## Adding a new service method

When the upstream UTxO RPC spec adds a method:

1. Wait for `go-codegen` to be regenerated and bumped in `go.mod`.
2. Add the wrapper in the appropriate `query.go` / `submit.go` / `sync.go` / `watch.go`.
3. Provide both `Method()` and `MethodWithContext(ctx)` forms.
4. Inject headers with `AddHeadersToRequest()`.
5. Mirror the same wrapper in `v1alpha/<file>.go` only if the upstream v1alpha spec also has it; otherwise leave v1alpha alone.
6. If the method is Cardano-relevant, add a high-level helper in `cardano/`.
7. Add an example under `examples/<service>/main.go` only if it demonstrates a new capability not already shown.

## Review checklist

For reviewers (human or agent) — what to look for in a PR:

- [ ] Conventional Commit title (CI will fail otherwise).
- [ ] `make test` passes locally; no new races.
- [ ] `make format` is clean — no gofumpt / gci diff against the PR.
- [ ] No new `golangci-lint` failures.
- [ ] If adding to root `sdk/`: both `Method()` and `MethodWithContext()` exist, headers are injected via `AddHeadersToRequest()`.
- [ ] If adding to `cardano/`: hex/base64 decoding is explicit and errors propagate; no silent `_ = ...` on decode results.
- [ ] Public API additions have GoDoc comments.
- [ ] No new features in `v1alpha/` — only bugfixes/back-compat.
- [ ] Breaking changes are marked with `!` in the commit and called out in the PR description with a migration note.
- [ ] No committed planning docs, scratch files, or `.claude/` artifacts.
- [ ] `go.mod` / `go.sum` changes are intentional and minimal (Dependabot or a real new dependency).
- [ ] No hard-coded API keys, URLs, or secrets in examples — read from env.

## Things to avoid

- Adding new features to `v1alpha/` — that tree is frozen.
- Editing generated code in `go-codegen` types — change upstream instead.
- Shipping examples that hard-code API keys or assume a specific endpoint.
- Adding a method only with `WithContext` or only without — always pair them.
- Skipping `AddHeadersToRequest()` in service wrappers — auth headers won't propagate and users will see opaque 401s.
- Introducing parallel HTTP clients or transports — go through `UtxorpcClient` so options compose.
- Using `fmt.Errorf` to wrap errors without `%w`, or comparing errors with `==`.
