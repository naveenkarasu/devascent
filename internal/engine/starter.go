package engine

import "devascent/internal/grader"

// GradingAvailable reports whether DevAscent has a function-call grader for a
// language (so the bench / Entrance Test / lessons can grade in it). Widens as
// per-language harnesses land.
func GradingAvailable(lang string) bool {
	switch lang {
	case "python", "go", "csharp", "javascript", "typescript", "java", "rust":
		return true
	}
	return false
}

// Starter generates the language-native starter stub for a code task when
// DevAscent generates one (e.g. Go infers a typed stub). ok is false for
// languages that keep their AUTHORED starter (python and reference-only langs) —
// callers must then leave the existing starter untouched rather than blank it.
// pySource (may be "") gives generated stubs nicer parameter names.
func Starter(lang, funcName, pySource string, tests []grader.TestCase, shape grader.Shape) (code string, ok bool) {
	switch lang {
	case "go":
		return grader.GoStarter(funcName, pySource, tests, shape), true
	case "csharp":
		return grader.CSharpStarter(funcName, pySource, tests, shape), true
	case "javascript":
		return grader.JSStarter(funcName, pySource, tests, shape), true
	case "typescript":
		return grader.TSStarter(funcName, pySource, tests, shape), true
	case "java":
		return grader.JavaStarter(funcName, pySource, tests, shape), true
	case "rust":
		return grader.RustStarter(funcName, pySource, tests, shape), true
	}
	return "", false
}
