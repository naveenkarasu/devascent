package content

import (
	"sort"
	"strconv"
	"strings"
)

// installLangOrder is the canonical display order for the install guide (matches
// the language-pick order). Any guide not listed here is appended alphabetically.
var installLangOrder = []string{"python", "javascript", "typescript", "go", "java", "csharp", "cpp", "rust"}

var installOSOrder = []struct{ key, label string }{
	{"windows", "Windows"},
	{"macos", "macOS"},
	{"linux", "Linux"},
}

// RenderInstallMarkdown renders all install guides to one Markdown document — the
// source for the repo INSTALL.md (ADR-0007: the file and the in-game Install Help
// screen come from the SAME data so they never drift). Deterministic order keeps
// the sync test stable.
func (c Catalog) RenderInstallMarkdown() string {
	byLang := map[string]InstallGuide{}
	for _, g := range c.InstallGuides {
		byLang[g.Lang] = g
	}
	order := append([]string{}, installLangOrder...)
	var extra []string
	for _, g := range c.InstallGuides {
		known := false
		for _, l := range installLangOrder {
			if l == g.Lang {
				known = true
				break
			}
		}
		if !known {
			extra = append(extra, g.Lang)
		}
	}
	sort.Strings(extra)
	order = append(order, extra...)

	var b strings.Builder
	b.WriteString("# Installing language toolchains for DevAscent\n\n")
	b.WriteString("DevAscent ships with **no bundled language runtimes** — it stays small and platform-independent, ")
	b.WriteString("and you install only the language(s) you want to play. When you pick a language, DevAscent checks ")
	b.WriteString("whether its toolchain is installed and actually works; if not, it points you here.\n\n")
	b.WriteString("You can read every primer and lesson without installing anything — a toolchain is only needed to ")
	b.WriteString("*run and grade* your own code in that language.\n\n")
	b.WriteString("<!-- Generated from internal/content/data/install/*.yaml — edit those, not this file. -->\n\n")

	for _, lang := range order {
		g, ok := byLang[lang]
		if !ok {
			continue
		}
		b.WriteString("## " + g.Label + "\n\n")
		if g.Notes != "" {
			b.WriteString("> " + g.Notes + "\n\n")
		}
		for _, o := range installOSOrder {
			st, ok := g.OS[o.key]
			if !ok {
				continue
			}
			b.WriteString("### " + o.label + "\n\n")
			if st.Link != "" {
				b.WriteString("Download: " + st.Link + "\n\n")
			}
			for i, s := range st.Steps {
				b.WriteString(strconv.Itoa(i+1) + ". " + s + "\n")
			}
			if st.Verify != "" {
				b.WriteString("\nVerify: `" + st.Verify + "`\n")
			}
			b.WriteString("\n")
		}
	}
	return b.String()
}
