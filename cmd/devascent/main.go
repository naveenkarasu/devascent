package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"devascent/internal/content"
	"devascent/internal/toolchain"
	"devascent/internal/tui"
)

// Build metadata, overridden at release time via -ldflags (see .goreleaser.yaml).
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// -primer <lang>: dump every primer for a language (syntax-highlighted) to
	// stdout and exit — a QA aid to eyeball a language without playing the game.
	primerLang := flag.String("primer", "", "preview all primers for a language (python|java|cpp|csharp|javascript|typescript|rust|go) and exit")
	advancedLang := flag.String("advanced", "", "preview Stage-2 Advanced Topics for a language and exit")
	genInstall := flag.Bool("geninstall", false, "regenerate INSTALL.md from internal/content/data/install/*.yaml and exit")
	doctor := flag.Bool("doctor", false, "probe all language toolchains (real compile+run) and print availability, then exit")
	showVersion := flag.Bool("version", false, "print version information and exit")
	flag.Parse()
	if *showVersion {
		fmt.Printf("devascent %s (commit %s, built %s)\n", version, commit, date)
		return
	}
	if *doctor {
		runDoctor()
		return
	}
	if *genInstall {
		cat, err := content.Load()
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		if err := os.WriteFile("INSTALL.md", []byte(cat.RenderInstallMarkdown()), 0o644); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		fmt.Println("wrote INSTALL.md")
		return
	}
	if *primerLang != "" {
		out, err := tui.PreviewPrimers(*primerLang)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		fmt.Print(out)
		return
	}
	if *advancedLang != "" {
		out, err := tui.PreviewAdvancedTopics(*advancedLang)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		fmt.Print(out)
		return
	}

	// AltScreen: render as a clean full-screen app instead of printing each
	// frame into scrollback (which made successive screens stack).
	if _, err := tea.NewProgram(tui.New(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

// runDoctor probes every supported language toolchain with the real compile+run
// canary and prints a report — the same detection the language picker uses.
func runDoctor() {
	det := toolchain.New()
	fmt.Println("DevAscent toolchain check (real compile + run probe):")
	fmt.Println()
	icon := map[toolchain.Status]string{
		toolchain.Available: "OK  ",
		toolchain.Missing:   "MISS",
		toolchain.Broken:    "WARN",
		toolchain.Unknown:   "??  ",
	}
	for _, lang := range det.Languages() {
		ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
		p := det.Capability(ctx, lang)
		cancel()
		line := fmt.Sprintf("  [%s] %-11s", icon[p.Status], lang)
		if p.Version != "" {
			line += "  " + p.Version
		}
		fmt.Println(line)
		if p.Status != toolchain.Available && p.Reason != "" {
			fmt.Printf("         %s\n", p.Reason)
		}
	}
	fmt.Println()
	fmt.Println("OK = installed & working · MISS = not found · WARN = found but the canary failed")
}
