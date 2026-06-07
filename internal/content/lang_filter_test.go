package content

import "testing"

// TestDiagnosticsForLang_MachineErrorIsLanguageSpecific verifies the per-slot
// language filter: machine-error serves the session language's native variants
// (a Go compile error for Go, a Rust panic for Rust), while language-neutral
// slots (coding-floor, spec, machine-terminal) fall through to their Python
// definitions for every language.
func TestDiagnosticsForLang_MachineErrorIsLanguageSpecific(t *testing.T) {
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}

	idsFor := func(lang string) map[string]bool {
		got := map[string]bool{}
		for _, d := range c.DiagnosticsForLang(lang) {
			got[d.ID] = true
		}
		return got
	}

	goIDs := idsFor("go")
	if !goIDs["go_err_index"] {
		t.Errorf("go intake should include the Go machine-error variant go_err_index")
	}
	if goIDs["err_index"] {
		t.Errorf("go intake should NOT include the Python machine-error variant err_index")
	}
	// language-neutral slot falls through to Python for Go
	if !goIDs["add"] {
		t.Errorf("go intake should still include the language-neutral coding-floor item add (Python fallback)")
	}

	pyIDs := idsFor("python")
	if !pyIDs["err_index"] {
		t.Errorf("python intake should include the Python machine-error variant err_index")
	}
	if pyIDs["go_err_index"] {
		t.Errorf("python intake should NOT include the Go machine-error variant go_err_index")
	}

	rustIDs := idsFor("rust")
	if !rustIDs["rust_err_move"] {
		t.Errorf("rust intake should include the Rust machine-error variant rust_err_move")
	}

	// cpp has no authored machine-error variants → falls back to the Python ones.
	cppIDs := idsFor("cpp")
	if !cppIDs["err_index"] {
		t.Errorf("cpp intake should fall back to the Python machine-error variant err_index")
	}
	if cppIDs["go_err_index"] {
		t.Errorf("cpp intake should not pick up another language's machine-error variant")
	}
}

// TestLessonsForLang_FallbackAndCoverage verifies every supported language yields
// a full 10-lesson Tutorial Island (its own variants where authored, Python
// fallback otherwise), in Order.
func TestLessonsForLang_FallbackAndCoverage(t *testing.T) {
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	for _, lang := range []string{"python", "go", "rust", "java", "csharp", "javascript", "typescript"} {
		ls := c.LessonsForLang(lang)
		if len(ls) != 10 {
			t.Fatalf("%s: expected 10 lessons, got %d", lang, len(ls))
		}
		want := normLang(lang)
		for i, l := range ls {
			if normLang(l.Lang) != want {
				t.Errorf("%s lesson %q resolved to lang %q, want %q", lang, l.ID, normLang(l.Lang), want)
			}
			if i > 0 && ls[i-1].Order > l.Order {
				t.Errorf("%s lessons out of Order at index %d", lang, i)
			}
		}
	}

	// An unauthored language (cpp) degrades to the full Python tutorial.
	cpp := c.LessonsForLang("cpp")
	if len(cpp) != 10 {
		t.Fatalf("cpp should fall back to 10 Python lessons, got %d", len(cpp))
	}
	for _, l := range cpp {
		if normLang(l.Lang) != "python" {
			t.Errorf("cpp fallback lesson %q should be Python, got %q", l.ID, l.Lang)
		}
	}
}
