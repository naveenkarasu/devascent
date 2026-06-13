package grader

import (
	"strings"
	"testing"
)

// ParseHarnessOutput must surface the player's own print/log output (every line
// except the result-marker line) so they can debug, while still grading from
// the marker.
func TestParseHarnessOutput_SurfacesPlayerStdout(t *testing.T) {
	tests := []TestCase{{Name: "case-1", Expected: float64(3)}}
	// Player printed two debug lines; the harness appended the marker line.
	out := "debug: a=1\ndebug: b=2\n" + marker + `[{"name":"case-1","got":3}]` + "\n"
	v := ParseHarnessOutput(out, tests)

	if !v.Passed {
		t.Fatalf("expected pass, got %+v", v)
	}
	if !strings.Contains(v.Stdout, "debug: a=1") || !strings.Contains(v.Stdout, "debug: b=2") {
		t.Fatalf("player output not surfaced: %q", v.Stdout)
	}
	if strings.Contains(v.Stdout, marker) {
		t.Fatalf("marker leaked into player output: %q", v.Stdout)
	}
}

func TestParseHarnessOutput_NoStdoutWhenSilent(t *testing.T) {
	tests := []TestCase{{Name: "case-1", Expected: float64(3)}}
	out := marker + `[{"name":"case-1","got":3}]`
	v := ParseHarnessOutput(out, tests)
	if v.Stdout != "" {
		t.Fatalf("expected empty stdout for a silent program, got %q", v.Stdout)
	}
}
