package grader

import "testing"

// Locks the line-ending normalization in stdout grading: a correct multi-line
// program must pass regardless of the platform's line separator (Java/C# emit
// \r\n on Windows; the authored signal uses \n). Without this, real Windows
// players — and the Windows native gate — fail correct solutions. Pure string
// logic, no toolchain.
func TestStdoutVerdict_NormalizesLineEndings(t *testing.T) {
	cases := []struct {
		name       string
		gotStdout  string
		want       string
		wantPassed bool
	}{
		{"crlf stdout vs lf signal", "a\r\nb\r\n", "a\nb", true},
		{"lone cr", "a\rb", "a\nb", true},
		{"trailing whitespace trimmed", "a\nb\n\n  ", "a\nb", true},
		{"genuine mismatch still fails", "a\r\nc", "a\nb", false},
		{"single line crlf", "42\r\n", "42", true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			v := stdoutVerdict(execOut{stdout: c.gotStdout, exit: 0}, c.want)
			if v.Passed != c.wantPassed {
				t.Errorf("stdoutVerdict(%q, %q).Passed = %v, want %v", c.gotStdout, c.want, v.Passed, c.wantPassed)
			}
		})
	}
}

func TestNormalizeNL(t *testing.T) {
	if got := normalizeNL("x\r\ny\r\n"); got != "x\ny" {
		t.Errorf("normalizeNL crlf = %q, want %q", got, "x\ny")
	}
	if got := normalizeNL("  trim\r\n  "); got != "trim" {
		t.Errorf("normalizeNL trim = %q, want %q", got, "trim")
	}
}
