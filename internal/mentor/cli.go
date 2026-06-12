package mentor

// CLI backends: the player's own installed AI agents, invoked exactly the way
// each vendor documents for scripting — unmodified binary, its own login, no
// credential files touched. The context pack travels over STDIN (never argv:
// Windows .cmd shims can't safely quote multi-line arguments), fixed flags
// keep every call tool-less/read-only, and the working directory is a neutral
// temp dir so no project config or CLAUDE.md leaks into the call.

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// runCLI executes one CLI call with the pack on stdin.
func runCLI(ctx context.Context, bin string, args []string, stdin string) (string, string, error) {
	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Dir = os.TempDir()
	cmd.Stdin = strings.NewReader(stdin)
	var out, errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	err := cmd.Run()
	return out.String(), errb.String(), err
}

// --- Claude Code -----------------------------------------------------------

type claudeBackend struct{ model string }

func newClaude(cfg Config) *claudeBackend {
	m := cfg.Model
	if m == "" {
		m = "haiku" // cheapest/fastest; hint calls are tiny
	}
	return &claudeBackend{model: m}
}

func (b *claudeBackend) ID() string   { return "claude" }
func (b *claudeBackend) Name() string { return "Claude Code (subscription)" }

func (b *claudeBackend) Present() (bool, string) {
	p, err := exec.LookPath("claude")
	if err != nil {
		return false, ""
	}
	return true, p
}

func (b *claudeBackend) Ask(ctx context.Context, prompt string) (string, error) {
	out, stderr, err := runCLI(ctx, "claude",
		[]string{"-p", "Respond to the request in the piped input.", "--output-format", "json", "--model", b.model},
		prompt)
	if err != nil {
		return "", fmt.Errorf("claude: %w: %s", err, firstLine(stderr))
	}
	var res struct {
		Result string `json:"result"`
	}
	if jerr := json.Unmarshal([]byte(lastJSONLine(out)), &res); jerr != nil || res.Result == "" {
		return "", fmt.Errorf("claude: unparseable output")
	}
	return res.Result, nil
}

// --- OpenAI Codex CLI ------------------------------------------------------

type codexBackend struct{ model string }

func newCodex(cfg Config) *codexBackend { return &codexBackend{model: cfg.Model} }

func (b *codexBackend) ID() string   { return "codex" }
func (b *codexBackend) Name() string { return "Codex CLI (ChatGPT account)" }

func (b *codexBackend) Present() (bool, string) {
	p, err := exec.LookPath("codex")
	if err != nil {
		return false, ""
	}
	return true, p
}

func (b *codexBackend) Ask(ctx context.Context, prompt string) (string, error) {
	tmp, err := os.CreateTemp("", "devascent-mentor-*.txt")
	if err != nil {
		return "", err
	}
	tmpPath := tmp.Name()
	tmp.Close()
	defer os.Remove(tmpPath)

	args := []string{"exec", "--sandbox", "read-only", "--skip-git-repo-check", "--output-last-message", tmpPath}
	if b.model != "" {
		args = append(args, "-m", b.model)
	}
	args = append(args, "-") // prompt from stdin
	_, stderr, err := runCLI(ctx, "codex", args, prompt)
	if err != nil {
		return "", fmt.Errorf("codex: %w: %s", err, firstLine(stderr))
	}
	data, err := os.ReadFile(tmpPath)
	if err != nil || len(bytes.TrimSpace(data)) == 0 {
		return "", fmt.Errorf("codex: no final message")
	}
	return string(data), nil
}

// --- GitHub Copilot CLI ----------------------------------------------------

type copilotBackend struct{ model string }

func newCopilot(cfg Config) *copilotBackend { return &copilotBackend{model: cfg.Model} }

func (b *copilotBackend) ID() string   { return "copilot" }
func (b *copilotBackend) Name() string { return "GitHub Copilot CLI" }

func (b *copilotBackend) Present() (bool, string) {
	p, err := exec.LookPath("copilot")
	if err != nil {
		return false, ""
	}
	return true, p
}

func (b *copilotBackend) Ask(ctx context.Context, prompt string) (string, error) {
	// Programmatic mode default-denies tools; -s strips session noise. No
	// --allow-tool flags are ever passed, so the call stays pure Q&A.
	args := []string{"-p", "Respond to the request in the piped input.", "-s"}
	if b.model != "" {
		args = append(args, "--model", b.model)
	}
	out, stderr, err := runCLI(ctx, "copilot", args, prompt)
	if err != nil {
		return "", fmt.Errorf("copilot: %w: %s", err, firstLine(stderr))
	}
	if strings.TrimSpace(out) == "" {
		return "", fmt.Errorf("copilot: empty output")
	}
	return out, nil
}

// --- helpers ---------------------------------------------------------------

func firstLine(s string) string {
	s = strings.TrimSpace(s)
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		s = s[:i]
	}
	if len(s) > 200 {
		s = s[:200]
	}
	return s
}

// lastJSONLine finds the final JSON object line in mixed output (CLIs may
// print banners before the payload).
func lastJSONLine(out string) string {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		t := strings.TrimSpace(lines[i])
		if strings.HasPrefix(t, "{") {
			return t
		}
	}
	return strings.TrimSpace(out)
}
