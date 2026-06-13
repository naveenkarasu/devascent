package content

import (
	"os"
	"testing"
)

// TestInstallGuidesLoad asserts every offered general-purpose language has a
// complete install guide (all three OSes, with a link, steps, and a verify
// command) — so a player who lands on an unavailable language always gets help.
func TestInstallGuidesLoad(t *testing.T) {
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	wantLangs := []string{"python", "javascript", "typescript", "go", "java", "csharp", "cpp", "rust"}
	for _, lang := range wantLangs {
		g, ok := c.InstallGuideForLang(lang)
		if !ok {
			t.Fatalf("missing install guide for %q", lang)
		}
		if g.Label == "" {
			t.Fatalf("%s: empty label", lang)
		}
		for _, osKey := range []string{"windows", "macos", "linux"} {
			st, ok := g.OS[osKey]
			if !ok {
				t.Fatalf("%s: missing %s instructions", lang, osKey)
			}
			if st.Link == "" {
				t.Fatalf("%s/%s: missing link", lang, osKey)
			}
			if len(st.Steps) == 0 {
				t.Fatalf("%s/%s: no steps", lang, osKey)
			}
			if st.Verify == "" {
				t.Fatalf("%s/%s: missing verify command", lang, osKey)
			}
		}
	}
}

// TestInstallMarkdownInSync guarantees the repo INSTALL.md is the rendered output
// of the install-guide data (single source of truth — edit the YAML, regenerate).
func TestInstallMarkdownInSync(t *testing.T) {
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	want := c.RenderInstallMarkdown()
	got, err := os.ReadFile("../../INSTALL.md")
	if err != nil {
		t.Fatalf("read INSTALL.md: %v (regenerate it from the install YAMLs)", err)
	}
	if string(got) != want {
		t.Fatalf("INSTALL.md is out of sync with the install YAMLs — regenerate it (see cmd/devascent -geninstall)")
	}
}
