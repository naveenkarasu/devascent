package content

import (
	"testing"
)

// TestPythonPrimersNoSemicolons guards an idiomatic-Python rule: a primer that
// teaches beginners must never use ';' to chain statements (PEP 8). Semicolons
// inside string/char literals (e.g. ";".join(x) or s.split(";")) are fine, so we
// only flag a ';' that appears OUTSIDE quotes.
func TestPythonPrimersNoSemicolons(t *testing.T) {
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range c.Primers {
		if p.Lang != "python" {
			continue
		}
		for _, sec := range p.Sections {
			for _, op := range sec.Ops {
				if semicolonOutsideQuotes(op.Code) {
					t.Errorf("python primer %q op %q uses a ';' statement separator (not idiomatic Python): %q",
						p.Category, op.Label, op.Code)
				}
			}
		}
	}
}

// semicolonOutsideQuotes reports whether s contains a ';' that is not inside a
// single- or double-quoted string literal.
func semicolonOutsideQuotes(s string) bool {
	var quote rune // 0 when not inside a string
	for _, r := range s {
		switch {
		case quote != 0:
			if r == quote {
				quote = 0
			}
		case r == '\'' || r == '"':
			quote = r
		case r == ';':
			return true
		}
	}
	return false
}
