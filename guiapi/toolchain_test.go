package guiapi

import (
	"os/exec"
	"testing"
)

// The gating surface: presence sweep covers every graded language with a valid
// status, every graded language has an install guide with steps for this OS,
// and the capability canary verifies a toolchain that is actually present.
func TestCapabilityGatingSurface(t *testing.T) {
	t.Setenv("DEVASCENT_SAVE_DIR", t.TempDir()) // never touch the real save
	e, err := New()
	if err != nil {
		t.Fatal(err)
	}

	statuses := e.LangPresence()
	if len(statuses) != len(GradedLanguages()) {
		t.Fatalf("presence sweep covered %d langs, want %d", len(statuses), len(GradedLanguages()))
	}
	valid := map[string]bool{"available": true, "missing": true, "broken": true, "unknown": true}
	for _, s := range statuses {
		if !valid[s.Status] {
			t.Errorf("%s: invalid status %q", s.Lang, s.Status)
		}
		if s.Verified {
			t.Errorf("%s: presence-only probe claims capability verification", s.Lang)
		}
	}

	for _, lang := range GradedLanguages() {
		g := e.InstallGuideFor(lang)
		if !g.Found || len(g.Steps) == 0 {
			t.Errorf("install guide missing/empty for %s on %s", lang, g.OS)
		}
	}

	if _, err := exec.LookPath("python"); err != nil {
		t.Skip("python not on PATH; skipping the capability round-trip")
	}
	st := e.CheckLang("python")
	if st.Status != "available" || !st.Verified {
		t.Fatalf("python capability check = %+v, want verified available", st)
	}
	// Recheck invalidates and still lands available.
	if st2 := e.RecheckLang("python"); st2.Status != "available" || !st2.Verified {
		t.Fatalf("python recheck = %+v", st2)
	}
}
