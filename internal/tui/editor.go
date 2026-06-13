package tui

import (
	"os"
	"os/exec"
	"runtime"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// editorOpt is one selectable editor in the in-game picker.
type editorOpt struct {
	label string // shown to the player
	cmd   string // command line (may include args, e.g. "code -w")
}

var editorOpts = []editorOpt{
	{"VS Code", "code -w"},
	{"nano", "nano"},
	{"vim", "vim"},
	{"notepad", "notepad"},
}

// editorAvailable reports whether the given option's executable is on PATH.
func editorAvailable(o editorOpt) bool {
	fields := strings.Fields(o.cmd)
	if len(fields) == 0 {
		return false
	}
	_, err := exec.LookPath(fields[0])
	return err == nil
}

// resolveEditor returns the command line to use: the explicit choice if set,
// else $VISUAL/$EDITOR, else a per-OS default.
func resolveEditor(choice string) string {
	if strings.TrimSpace(choice) != "" {
		return choice
	}
	if ed := os.Getenv("VISUAL"); ed != "" {
		return ed
	}
	if ed := os.Getenv("EDITOR"); ed != "" {
		return ed
	}
	if runtime.GOOS == "windows" {
		return "notepad"
	}
	return "vi"
}

// markerBody is the stable text of the line separating the (read-only) task
// header from the player's code. The comment prefix varies by language, so
// unwrapFromEditor matches on this body (prefix-agnostic).
const markerBody = "===== write your solution below this line ====="

// editorMarker is the Python-style marker line (kept for reference/tests).
const editorMarker = "# " + markerBody

// commentPrefix returns the line-comment token for a source extension: "#" for
// Python, "//" for every other supported language (Java/Go/C#/JS/TS/Rust/C++).
func commentPrefix(ext string) string {
	if ext == ".py" {
		return "#"
	}
	return "//"
}

// wrapForEditor prepends the task prompt as a commented header so the player can
// see the question WHILE editing. The header uses the language's comment style
// (prefix), and is stripped on save — so it never reaches the compiler either way.
func wrapForEditor(initial, prompt, prefix string) string {
	if prefix == "" {
		prefix = "#"
	}
	var b strings.Builder
	b.WriteString(prefix + " +-- TASK ----------------------------------------------------\n")
	for _, ln := range strings.Split(strings.TrimRight(prompt, "\n"), "\n") {
		b.WriteString(prefix + " | " + ln + "\n")
	}
	b.WriteString(prefix + " +-----------------------------------------------------------\n")
	b.WriteString(prefix + " " + markerBody + "\n")
	b.WriteString(initial)
	return b.String()
}

// unwrapFromEditor strips the task header back off on save (so it never gets
// graded or re-accumulated), matching the marker body regardless of comment
// prefix. If the marker is gone, return the file as-is.
func unwrapFromEditor(content string) string {
	if i := strings.LastIndex(content, markerBody); i >= 0 {
		rest := content[i+len(markerBody):]
		if nl := strings.IndexByte(rest, '\n'); nl >= 0 {
			return rest[nl+1:]
		}
		return ""
	}
	return content
}

// editorCmd writes the prompt+code to a temp file (with the given extension, so
// the editor highlights it correctly), opens the chosen editor, and on exit reads
// the (unwrapped) code back as an editorFinishedMsg. In tests this path is
// bypassed by sending editorFinishedMsg directly. ext includes the dot (".py").
func editorCmd(initial, choice, prompt, ext string) tea.Cmd {
	if ext == "" {
		ext = ".py"
	}
	f, err := os.CreateTemp("", "devascent-*"+ext)
	if err != nil {
		return func() tea.Msg { return editorFinishedMsg{err: err} }
	}
	path := f.Name()
	_, _ = f.WriteString(wrapForEditor(initial, prompt, commentPrefix(ext)))
	_ = f.Close()

	ed := resolveEditor(choice)
	// The command may include args (e.g. "code -w"). Split into exe + args,
	// append the file path last. Resolve the exe via PATH so Windows shims
	// (e.g. code.cmd) are found.
	fields := strings.Fields(ed)
	exe := fields[0]
	if resolved, lpErr := exec.LookPath(exe); lpErr == nil {
		exe = resolved
	}
	args := append(fields[1:], path)

	c := exec.Command(exe, args...)
	return tea.ExecProcess(c, func(runErr error) tea.Msg {
		data, readErr := os.ReadFile(path)
		_ = os.Remove(path)
		if runErr == nil {
			runErr = readErr
		}
		return editorFinishedMsg{code: unwrapFromEditor(string(data)), err: runErr}
	})
}

// editorHint returns a one-line note on the current editor + how to change it.
func editorHint(choice string) string {
	if strings.TrimSpace(choice) != "" {
		return "(editor: " + choice + " — [c] to change)"
	}
	return "([c] to pick an editor)"
}
