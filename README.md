# DevAscent

[![CI](https://github.com/naveenkarasu/devascent/actions/workflows/ci.yml/badge.svg)](https://github.com/naveenkarasu/devascent/actions/workflows/ci.yml)
[![License: AGPL v3](https://img.shields.io/badge/License-AGPL_v3-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.26%2B-00ADD8?logo=go&logoColor=white)](go.mod)

**A terminal-based software-developer career simulator.** Start as a junior and
grow toward senior/staff by doing the *real* work — writing code, solving problems,
reading briefs, working the terminal — all graded honestly by a real compiler, not
by a quiz. No abstraction, no fake editor: it's you, a terminal, and real tests.

> **Status: early.** The core loop (orientation → DSA apprenticeship) is playable
> and graded in 7 languages. Prebuilt downloads aren't published yet — for now you
> build it from source (one command, below). Expect rough edges.

---

## What it is

DevAscent runs entirely in your terminal (a TUI built with [Bubble Tea](https://github.com/charmbracelet/bubbletea)).
A run flows through:

1. **The Entrance Test** — a short adaptive diagnostic that figures out where you are.
2. **Tutorial Island** (for absolute beginners) or **Dev-Literacy** (terminal/git
   basics) — shown only if you need them.
3. **The Apprenticeship** — a curated DSA bench (Blind 75 / NeetCode 150-style)
   where you write real functions and they're graded by real, hidden tests.

You write your solutions in **your own editor**, and they're compiled and run by
**your own language toolchain** — so the grading is exactly as honest as a real job.

## Pick your language

DevAscent ships **no language runtimes** — the download stays tiny, and you install
the toolchain only for the language you want to play in. When you choose a language,
DevAscent checks whether its toolchain is installed and actually works; if not, it
shows you how to install it. You can read every primer and lesson **without
installing anything** — a toolchain is only needed to *run and grade* your code.

**Languages with full graded play:** Python · Go · Rust · Java · C# · JavaScript · TypeScript

C++ currently has reference content (primers and topics you can read) but graded
play isn't supported yet.

See [`INSTALL.md`](INSTALL.md) for per-OS install steps. Once you're in the game,
choosing a language you haven't installed opens its install guide automatically.

You'll also want a text editor. DevAscent opens your `$EDITOR` / `$VISUAL` (e.g.
`code -w` for VS Code) and has an in-game editor picker if you haven't set one.

## Two ways to play

DevAscent ships **two frontends over the same game** — same engine, same grading,
same save files. Progress carries between them (one save slot per language), so
pick whichever fits, or switch freely:

| | `devascent` (terminal) | `devascent-gui` (desktop) |
|---|---|---|
| Feel | TUI in your terminal, code in your own `$EDITOR` | Desktop app with a built-in Monaco code editor |
| Platforms | Windows / macOS / Linux · amd64 + arm64 | Windows (amd64) · macOS (universal) · Linux (amd64) |
| Footprint | One static binary, zero dependencies | Native webview (WebView2 / WebKitGTK / Cocoa) |

## Get started

### Option A — download a prebuilt build (easiest)

Grab what you want from the [**Releases**](https://github.com/naveenkarasu/devascent/releases/latest)
page — `devascent_*` archives are the terminal app, `devascent-gui_*` archives
are the desktop app (each release carries both; verify against the checksums).

**Terminal app:**

```sh
# macOS / Linux
tar -xzf devascent_*_*.tar.gz
./devascent -doctor   # see which language toolchains you have
./devascent
```

**Desktop app:** unzip and run `devascent-gui.exe` (Windows), `devascent-gui.app`
(macOS), or `devascent-gui` (Linux).

On **macOS** both builds are unsigned — run `xattr -d com.apple.quarantine <file>`
once, or right-click → Open. On **Windows**, SmartScreen → "More info" →
"Run anyway".

### Option B — build from source

**Terminal app** — needs [Go](https://go.dev/dl/) **1.26+** (no CGO; content is
embedded in the binary):

```sh
git clone https://github.com/naveenkarasu/devascent
cd devascent
go build -o devascent ./cmd/devascent
./devascent
```

**Desktop app** — additionally needs [Node 20+](https://nodejs.org) and the
[Wails CLI](https://wails.io) (Linux also: `libgtk-3-dev libwebkit2gtk-4.0-dev`):

```sh
go install github.com/wailsapp/wails/v2/cmd/wails@v2.12.0
cd devascent/devascent-gui
wails build        # → build/bin/devascent-gui
```

> `brew` / `scoop` install are planned next.

## A few flags (terminal app)

| Flag | What it does |
|---|---|
| `-doctor` | Check every language toolchain on your machine (real compile + run) and print what's installed and working. |
| `-primer <lang>` | Print a language's reference primers to the terminal and exit (peek without playing). |

Your progress is saved automatically in your per-OS config directory — one save
slot per language, shared by both frontends; set `DEVASCENT_SAVE_DIR` to
override where.

## License

[GNU AGPL-3.0](LICENSE). In short: it's free and open; if you distribute a modified
version — including running it as a network service — you must offer users the
corresponding source. DSA problem statements and lesson content are originally
authored for DevAscent.

## Support

DevAscent is free and built by one developer. If it helps you, support is welcome —
*(donation link coming soon: Ko-fi / GitHub Sponsors).*
