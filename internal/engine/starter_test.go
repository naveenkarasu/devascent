package engine

import (
	"strings"
	"testing"

	"devascent/internal/grader"
)

// Golden-ish guard for the language→starter dispatch. This is the cheap net the
// native compilation gate doesn't give us here: it pins the ok-semantics that
// keep python/reference langs from having their authored starter blanked, and
// confirms each generated language actually produces a stub. Pure string
// generation — no toolchain, so it's fast.

func TestGradingAvailable(t *testing.T) {
	for _, lang := range []string{"python", "go", "csharp", "javascript", "typescript", "java", "rust"} {
		if !GradingAvailable(lang) {
			t.Errorf("GradingAvailable(%q) = false, want true", lang)
		}
	}
	for _, lang := range []string{"cpp", "ruby", "", "Go", "kotlin"} {
		if GradingAvailable(lang) {
			t.Errorf("GradingAvailable(%q) = true, want false", lang)
		}
	}
}

func TestStarterGeneratedLangs(t *testing.T) {
	var shape grader.Shape // zero shape = plain problem
	for _, lang := range []string{"go", "csharp", "javascript", "typescript", "java", "rust", "cpp"} {
		code, ok := Starter(lang, "twoSum", "", nil, shape)
		if !ok {
			t.Errorf("Starter(%q) ok = false, want true", lang)
		}
		if strings.TrimSpace(code) == "" {
			t.Errorf("Starter(%q) produced empty stub", lang)
		}
	}
}

func TestStarterKeepsAuthoredStarter(t *testing.T) {
	// python and unknown languages must report ok=false so the caller leaves
	// the authored starter in place. C++ moved OUT of this set 2026-06-12:
	// reference-only still means the player sees C++ syntax, never Python.
	for _, lang := range []string{"python", "ruby", ""} {
		code, ok := Starter(lang, "twoSum", "", nil, grader.Shape{})
		if ok {
			t.Errorf("Starter(%q) ok = true, want false (authored starter must be kept)", lang)
		}
		if code != "" {
			t.Errorf("Starter(%q) code = %q, want empty", lang, code)
		}
	}
}
