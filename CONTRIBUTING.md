# Contributing to DevAscent

Thanks for your interest! DevAscent is an early, solo-built project, so the most
useful contributions right now are **bug reports, playtesting feedback, and small
focused fixes**. Larger changes are welcome too — please open an issue to discuss
direction before a big PR.

By contributing, you agree your contributions are licensed under the project's
[AGPL-3.0](LICENSE).

## Project layout

| Path | What's there |
|---|---|
| `cmd/devascent/` | Entry point + CLI flags (`-doctor`, `-primer`, …). |
| `internal/tui/` | Bubble Tea state machine + views, editor shell-out, content selection. |
| `internal/content/` | Catalog schema + the embedded `data/**` YAML (lessons, diagnostics, bench, primers, install guides). |
| `internal/grader/` | The `Grader` interface + backends. Per-language harnesses (Model A function-call grading). |
| `internal/toolchain/` | Detects the player's installed language toolchains (BYO model). |
| `internal/save/` | JSON save/resume. |

See [`ARCHITECTURE.md`](ARCHITECTURE.md) for how grading and toolchain detection work.

## Dev setup

You need [Go](https://go.dev/dl/) **1.26+**. No CGO; content is embedded via `go:embed`.

```sh
go build ./cmd/devascent     # build the binary
go run ./cmd/devascent       # run it
go run ./cmd/devascent -doctor   # see which language toolchains you have
```

DevAscent grades code with **your own toolchains**. The fast test suite only needs
Go + Python; the full suite compiles/runs real code in every language, so install
the toolchain for any language you touch (`-doctor` tells you what's missing).

## Before you open a PR

```sh
gofmt -l .                   # must print nothing
go vet ./...
go test ./... -short         # fast suite (Go + Python only)
go test ./...                # full suite — needed if you touched the grader or content
                             # (compiles/runs every language; slower)
```

- **Keep the grader gate green.** Tests in `internal/content` run reference
  solutions through the real toolchains and assert they pass — that's the safety
  net for any content or grader change.
- **Content conventions.** Per-language lesson variants keep their `tests:`
  byte-identical to the Go reference (the documented exceptions are the `dicts`
  and `read-the-crash` lessons for statically-typed languages). Every lesson task
  ships a `solution:` so the native round-trip test can validate it.
- **Match the surrounding style.** Small, focused commits with clear messages.

## Reporting bugs / proposing features

Use the issue templates (they ask for your OS, language toolchain, and `-doctor`
output — those make a bug actionable on a cross-platform tool). For anything
security-sensitive, see [`SECURITY.md`](SECURITY.md) instead of a public issue.
