package tui

import (
	"strings"
	"testing"
)

func TestEditorWrapUnwrap_RoundTrip(t *testing.T) {
	code := "def add(a, b):\n    return a + b\n"
	prompt := "Write a function add(a, b) that returns their sum."
	wrapped := wrapForEditor(code, prompt, "#")
	if !strings.Contains(wrapped, "add(a, b)") || !strings.Contains(wrapped, markerBody) {
		t.Fatalf("wrapped file should show the prompt + marker:\n%s", wrapped)
	}
	if got := unwrapFromEditor(wrapped); got != code {
		t.Fatalf("round-trip lost code:\n got %q\nwant %q", got, code)
	}
}

func TestEditorUnwrap_NoMarker_ReturnsAsIs(t *testing.T) {
	raw := "def f():\n    return 1\n"
	if got := unwrapFromEditor(raw); got != raw {
		t.Fatalf("no-marker content should pass through, got %q", got)
	}
}

func TestEditorUnwrap_MultilinePrompt(t *testing.T) {
	code := "def second_largest(nums):\n    return sorted(set(nums))[-2]\n"
	prompt := "Write second_largest(nums) that returns the\nsecond-largest DISTINCT value."
	if got := unwrapFromEditor(wrapForEditor(code, prompt, "#")); got != code {
		t.Fatalf("multiline-prompt round-trip failed:\n got %q", got)
	}
}

// Non-Python languages must use // comments (not #), and the header must still
// strip cleanly on unwrap so it never reaches the compiler.
func TestEditorWrap_PerLanguageComment(t *testing.T) {
	cases := map[string]string{".py": "#", ".java": "//", ".go": "//", ".cs": "//", ".js": "//", ".ts": "//", ".rs": "//", ".cpp": "//"}
	code := "class Solution { public long f(){ return 0L; } }\n"
	prompt := "Implement f()."
	for ext, prefix := range cases {
		wrapped := wrapForEditor(code, prompt, commentPrefix(ext))
		if !strings.Contains(wrapped, prefix+" +-- TASK") {
			t.Fatalf("%s: header should use %q comments:\n%s", ext, prefix, wrapped)
		}
		if ext != ".py" && strings.Contains(wrapped, "# +-- TASK") {
			t.Fatalf("%s: must NOT use Python # comments", ext)
		}
		if got := unwrapFromEditor(wrapped); got != code {
			t.Fatalf("%s: header must strip cleanly, got %q", ext, got)
		}
	}
}
