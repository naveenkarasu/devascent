package mentor

// Package mentor is the Track-A4 BYO-AI seam. The player's OWN AI answers
// tier-2/3 hints, write-up follow-ups, and review notes; the game ships no AI
// and holds no credentials — backends are the unmodified CLIs/servers the
// player already has (claude, codex, copilot, Ollama, any OpenAI-compatible
// endpoint), invoked tool-less/read-only. A built-in template backend keeps
// everything working offline, and tier-1 nudges are ALWAYS template-served.
//
// Binding guardrails (claudedocs/design_a2_hint_economy_a4_byo_ai_2026-06-12.md):
// the AI never grades, banks, or mutates saves — its output is display text
// only; the request fixes exactly what context is sent; responses are
// validated deterministically and fall back to templates on any violation,
// garbage, or timeout. The AI never sees hidden tests, canonical solutions,
// or the MCQ answer key — Request has no fields for them.

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"devascent/internal/save"
)

// Kind is what the mentor is being asked to do.
type Kind string

const (
	KindStrategy    Kind = "strategy"    // tier-2: approach, no code
	KindWalkthrough Kind = "walkthrough" // tier-3: steps + pseudocode, no compilable solution
	KindFollowup    Kind = "followup"    // A1: one probing question about the write-up
	KindReview      Kind = "review"      // post-bank: one strength + one improvement
)

// AskTimeout bounds every backend call; on expiry the caller falls back to
// templates (and refunds the token).
const AskTimeout = 45 * time.Second

// ProbeTimeout is longer: local models may need to load.
const ProbeTimeout = 90 * time.Second

// Canary is the capability-probe sentinel (same pattern as the toolchain
// compile-and-run canary).
const Canary = "DEVASCENT_OK"

// Request is the COMPLETE context a mentor call may see — adding a field here
// is a guardrails change and needs a design update first.
type Request struct {
	Kind       Kind
	Lang       string // player's language (for "no compilable solution" enforcement)
	Title      string
	Prompt     string // public problem statement
	Category   string
	Difficulty string
	PlayerCode string // the player's CURRENT editor code
	FailedRuns int    // failed grade attempts so far
	FirstFail  string // name of the first failing test case ("" if none) — never its data
	Writeup    string // player's write-up text (followup/review)
	Attempt    int    // nudge escalation counter (template tiering)
}

// Response is what the player sees.
type Response struct {
	Text     string `json:"text"`
	Source   string `json:"source"`   // backend ID or "template"
	FellBack bool   `json:"fellBack"` // an AI backend was selected but templates answered
}

// Backend is one way to reach an AI. Implementations must be read-only
// toward the player's machine (tool-less CLIs, sandboxed exec, plain HTTP).
type Backend interface {
	ID() string
	Name() string
	// Present is the cheap phase-1 check (binary on PATH / port answering).
	// info carries a version or address for the picker UI.
	Present() (present bool, info string)
	// Ask sends one prompt and returns the raw reply.
	Ask(ctx context.Context, prompt string) (string, error)
}

// Config is the machine-level mentor selection, stored as mentor.json next to
// the save slots (NOT per-language: the machine's AI doesn't change per run).
type Config struct {
	Backend  string `json:"backend"`            // backend ID; "" = templates only
	Endpoint string `json:"endpoint,omitempty"` // openai-compat base URL (e.g. http://localhost:1234/v1)
	Model    string `json:"model,omitempty"`    // model name (HTTP backends; overrides CLI defaults)
	APIKey   string `json:"api_key,omitempty"`  // openai-compat only; local servers ignore it. Player-entered, opt-in.
}

func configPath() (string, error) {
	d, err := save.Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(d, "mentor.json"), nil
}

// LoadConfig reads mentor.json; a missing file is the zero config.
func LoadConfig() (Config, error) {
	p, err := configPath()
	if err != nil {
		return Config{}, err
	}
	data, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, nil
		}
		return Config{}, err
	}
	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return Config{}, nil // corrupt config: fall back to templates, don't block the game
	}
	return c, nil
}

// SaveConfig writes mentor.json atomically (same discipline as save slots).
func SaveConfig(c Config) error {
	p, err := configPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	tmp := p + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, p)
}
