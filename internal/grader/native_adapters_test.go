package grader

import (
	"context"
	"strings"
	"testing"
	"time"

	"devascent/internal/toolchain"
)

func nativeGrader(t *testing.T) *LocalToolchain {
	t.Helper()
	g := NewLocalToolchain(toolchain.New())
	g.timeout = 120 * time.Second // real compiles (esp. a cold dotnet restore) need headroom
	return g
}

func skipUnlessAvailable(t *testing.T, g *LocalToolchain, lang string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	if g.det.Capability(ctx, lang).Status != toolchain.Available {
		t.Skipf("%s toolchain not available; skipping native grade test", lang)
	}
}

// TestNativeStdout proves each installed language can RUN player code and be
// graded by stdout — the universal oracle check (ADR-0007). Skips absent langs.
func TestNativeStdout(t *testing.T) {
	if testing.Short() {
		t.Skip("native multi-language grading is slow; skipped in -short")
	}
	src := map[string]string{
		"javascript": `console.log("hi")` + "\n",
		"typescript": `const m: string = "hi"` + "\n" + `console.log(m)` + "\n",
		"go":         "package main\nimport \"fmt\"\nfunc main() { fmt.Println(\"hi\") }\n",
		"java":       "public class Main { public static void main(String[] a) { System.out.println(\"hi\"); } }\n",
		"csharp":     "System.Console.WriteLine(\"hi\");\n",
		"rust":       "fn main() { println!(\"hi\"); }\n",
	}
	g := nativeGrader(t)
	for lang, source := range src {
		t.Run(lang, func(t *testing.T) {
			skipUnlessAvailable(t, g, lang)
			v, err := g.Grade(GradeRequest{Lang: lang, Check: CheckStdout, Source: source, Signal: "hi"})
			if err != nil {
				t.Fatal(err)
			}
			if !v.Passed {
				t.Fatalf("%s stdout grade should pass, got %+v (err=%q)", lang, v.Results, v.Err)
			}
		})
	}
}

// TestNativeCompiles proves the compile-check oracle: valid code compiles, broken
// code does not. Covers the fast compilers (Go/Rust/Java/TypeScript).
func TestNativeCompiles(t *testing.T) {
	if testing.Short() {
		t.Skip("native compile checks are slow; skipped in -short")
	}
	good := map[string]string{
		"go":         "package main\nfunc main() {}\n",
		"rust":       "fn main() {}\n",
		"java":       "public class Main { public static void main(String[] a) {} }\n",
		"typescript": "const x: number = 1\n",
	}
	bad := map[string]string{
		"go":         "package main\nfunc main() { x := }\n",     // syntax error
		"rust":       "fn main() { let x: i32 = \"s\"; }\n",      // type mismatch
		"java":       "public class Main { this is not java }\n", // garbage
		"typescript": "const x: number = \"not a number\"\n",     // type error
	}
	g := nativeGrader(t)
	for lang := range good {
		t.Run(lang, func(t *testing.T) {
			skipUnlessAvailable(t, g, lang)
			ok, _ := g.Grade(GradeRequest{Lang: lang, Check: CheckCompiles, Source: good[lang]})
			if !ok.Passed {
				t.Fatalf("%s: valid code should compile, got %+v", lang, ok.Results)
			}
			no, _ := g.Grade(GradeRequest{Lang: lang, Check: CheckCompiles, Source: bad[lang]})
			if no.Passed {
				t.Fatalf("%s: broken code should NOT compile", lang)
			}
		})
	}
}

// TestRustCompileError proves the compiler-as-oracle pattern (trybuild-style):
// broken Rust must fail with the EXPECTED diagnostic code (E0382, borrow of moved
// value) — this is exactly how the authored Rust advanced-topics exercises grade.
func TestRustCompileError(t *testing.T) {
	if testing.Short() {
		t.Skip("rust compile-error check is slow; skipped in -short")
	}
	g := nativeGrader(t)
	skipUnlessAvailable(t, g, "rust")
	moved := "fn main() {\n    let s = String::from(\"x\");\n    let t = s;\n    println!(\"{}\", s);\n    let _ = t;\n}\n"
	v, err := g.Grade(GradeRequest{Lang: "rust", Check: CheckCompileError, Source: moved, Signal: "E0382"})
	if err != nil {
		t.Fatal(err)
	}
	if !v.Passed {
		t.Fatalf("moved-value Rust should fail with E0382, got %+v", v.Results)
	}
	// And clean Rust should NOT satisfy a compile-error check.
	clean, _ := g.Grade(GradeRequest{Lang: "rust", Check: CheckCompileError, Source: "fn main() {}\n", Signal: "E0382"})
	if clean.Passed {
		t.Fatalf("clean Rust should not satisfy a compile-error check")
	}
	if !strings.Contains(strings.ToLower(clean.Results[0].Err), "compiled") {
		t.Logf("note: clean-compile rejection message: %q", clean.Results[0].Err)
	}
}
