package grader

import "strings"

// inlang.go — shared support for languages WITHOUT a stdlib JSON library
// (Rust, Java). Instead of serializing the result back to JSON, the harness
// embeds the EXPECTED value as a native literal, compares IN-LANGUAGE
// (==/equals/deepEquals), and prints one marker line per case:
//
//	<marker><name><US>("1"|"0")<US><gotRepr>
//
// where US is the unit separator (0x1f). ParseInLangOutput rebuilds the Verdict.
// This sidesteps hand-rolled JSON serialization (advisor's recommendation): native
// equality is trivially correct; the gotRepr is a debug string for failure display.

const inlangUS = "\x1f"

// ParseInLangOutput builds a Verdict from the in-language line protocol. If no
// marker line appears at all, the raw output is treated as a compile/runtime
// error (parity with ParseHarnessOutput's no-marker behavior).
func ParseInLangOutput(out string, tests []TestCase) Verdict {
	type rec struct {
		ok  bool
		got string
	}
	seen := map[string]rec{}
	found := false
	for _, ln := range strings.Split(out, "\n") {
		i := strings.Index(ln, marker)
		if i < 0 {
			continue
		}
		found = true
		rest := strings.TrimRight(ln[i+len(marker):], "\r")
		parts := strings.SplitN(rest, inlangUS, 3)
		if len(parts) < 2 {
			continue
		}
		r := rec{ok: parts[1] == "1"}
		if len(parts) == 3 {
			r.got = parts[2]
		}
		seen[parts[0]] = r
	}
	if !found {
		return Verdict{Err: trimErr(out)}
	}
	v := Verdict{Passed: true}
	for _, tc := range tests {
		cr := CaseResult{Name: tc.Name}
		if r, ok := seen[tc.Name]; ok {
			cr.Passed = r.ok
			cr.Got = r.got
		} else {
			cr.Passed = false
			cr.Err = "no result"
		}
		if !cr.Passed {
			v.Passed = false
		}
		v.Results = append(v.Results, cr)
	}
	if len(v.Results) == 0 {
		v.Passed = false
		v.Err = "no test results"
	}
	return v
}
