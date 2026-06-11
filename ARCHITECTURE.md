# Architecture

A map of how DevAscent fits together, for contributors. The core is a UI-agnostic
Go engine — an embedded content catalog and a grader that shells out to the
player's own language toolchains — with **two frontends over the same core**: a
terminal app (single static binary) and a Wails desktop app. They share the
grading, the content, and the save files.

## Packages

- **`internal/tui`** — the [Bubble Tea](https://github.com/charmbracelet/bubbletea)
  state machine and views. Drives the run (entrance test → tutorial/dev-literacy →
  bench), opens the player's `$EDITOR`, and selects content per session language.
- **`internal/content`** — the catalog schema and the `data/**` YAML, embedded with
  `go:embed`. Holds lessons, diagnostics, the DSA bench, primers, advanced topics,
  and the per-OS install guides. Most validation tests live here.
- **`internal/grader`** — the `Grader` interface and its backends, plus one harness
  per language.
- **`internal/toolchain`** — detects which language toolchains the player has.
- **`internal/save`** — JSON save/resume in the per-OS config dir; one slot per
  language, shared by both frontends.
- **`guiapi`** — the public facade the desktop app calls into (the GUI is a
  separate Go module, so it can't import `internal/*`). Wraps content + engine +
  grader + saves as stateful, JSON-friendly sessions.
- **`devascent-gui/`** — the [Wails](https://wails.io) desktop app: a Go backend
  binding `guiapi` to a Monaco-editor web frontend. Its own module
  (`replace devascent => ../`), so the core stays CGO-free.

## Grading (Model A — function-call, LeetCode-style)

The player writes a function; DevAscent grades by **calling it with each test's
arguments and comparing the result**. For a submission, the language harness
generates a complete throwaway program that:

1. embeds each test's arguments as **native literals** of an inferred type (no
   runtime input parsing — see the shared `gtype` inference from the test data),
2. defines the player's function,
3. calls it per test and emits a line-protocol result.

Two comparison strategies, by language:

- **JSON compare** (Python, Go, C#, JavaScript, TypeScript): the harness serializes
  the result as JSON; `ParseHarnessOutput` + `jsonEqual` compare it to the expected
  value with the same semantics across languages.
- **In-language compare** (Rust, Java): these have no stdlib JSON, so the harness
  embeds the *expected* value as a native literal too and compares in-language
  (`==` / `Arrays.equals` / `PartialEq`), emitting pass/fail via `ParseInLangOutput`.

Data-structure problems (`shape: linkedlist | tree`) inject a per-language prelude
that builds the node structure from an array and dumps it back for comparison.
Each test case is isolated so one runtime error doesn't sink the whole submission.

### Backends (`Grader` interface)

- **`LocalToolchain`** — the default. Dispatches to a per-language adapter that
  compiles/runs via the player's installed toolchain (the BYO model).
- **`NativePython`** — Python via the system interpreter; used by the fast test
  gate. `DEVASCENT_GRADER=native` selects it.
- **`WazeroPython`** — a sandboxed Python-in-WASM backend, kept opt-in
  (`DEVASCENT_GRADER=wazero`, needs a local `python.wasm`).

## Toolchain detection (BYO runtimes)

DevAscent ships no runtimes. `internal/toolchain` probes in two phases:

- **Presence** — fast: is the binary on `PATH`, and what version?
- **Capability** — deep: write a tiny canary program, actually compile **and** run
  it, and require it to print the `DEVASCENT_OK` sentinel. This catches a
  JRE-without-JDK, a missing linker, or a broken `PATH` *before* the player hits it
  mid-game. Results are cached; `-doctor` runs the same probe and prints a report.

Two independent axes keep concerns separate: **runtime availability** (is the
toolchain installed?) vs **grading maturity** (does a grader adapter exist?). A
language is playable-and-graded only when both hold; otherwise it's reference-only.

## Content & per-language variants

Lessons and the `machine-error` diagnostics have per-language variants selected by
the session language, with a **Python fallback** when a language has no authored
variant (`LessonsForLang` / `DiagnosticsForLang`). Primers are per-category and
per-language. The non-Python lesson tasks share one harness-safe shape, validated
end-to-end by the native round-trip tests that grade every reference `solution`
through the real toolchain.
