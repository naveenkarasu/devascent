package grader

import (
	"path/filepath"
	"runtime"
)

// nativeAdapter grades a language by shelling out to its real toolchain (ADR-0007)
// for the "compiler-as-oracle" checks: stdout (run + compare), compiles (build
// must succeed), and compile-error (build must fail with the expected diagnostic,
// the trybuild/compiletest pattern). It is driven by a per-language profile, so
// one implementation covers JS / TS / Go / Java / C# / Rust.
//
// Function-call (CheckTests) grading is NOT handled here: generating a typed
// driver that calls an arbitrary player function needs per-problem type metadata
// that statically-typed languages require and the bench doesn't carry yet. Python
// keeps CheckTests via pythonAdapter; everything else uses the oracle checks above.
type nativeAdapter struct {
	lang  string
	file  string                             // source filename written from req.Source
	extra func(dir string) map[string]string // optional aux files (e.g. a .csproj)
	// runCmds: compile (optional) then RUN — the last step's stdout is the output.
	runCmds func(r runner) [][]string
	// compCmds: compile ONLY — the last step's exit/stderr is judged.
	compCmds func(r runner) [][]string
	// testsDriver, if set, builds a full function-call grading program (player
	// source + harness) for Check==tests; the result is written to `file` and run
	// via runCmds, then parsed with ParseHarnessOutput (Model A, LeetCode-style).
	testsDriver func(req GradeRequest) (string, error)
	// inLang routes the tests output through ParseInLangOutput instead of
	// ParseHarnessOutput — for languages without stdlib JSON (Rust/Java) whose
	// harness compares in-language and prints the line protocol.
	inLang bool
}

func (a nativeAdapter) grade(req GradeRequest, r runner) (Verdict, error) {
	if err := r.write(a.file, req.Source); err != nil {
		return Verdict{}, err
	}
	if a.extra != nil {
		for name, content := range a.extra(r.dir) {
			if err := r.write(name, content); err != nil {
				return Verdict{}, err
			}
		}
	}
	switch req.Check {
	case CheckStdout:
		last, compileErr := execSteps(r, a.runCmds(r))
		if compileErr != "" {
			return Verdict{Err: compileErr}, nil
		}
		return stdoutVerdict(last, req.Signal), nil
	case CheckCompiles:
		last, _ := execSteps(r, a.compCmds(r))
		return compilesVerdict(last), nil
	case CheckCompileError:
		last, _ := execSteps(r, a.compCmds(r))
		return compileErrorVerdict(last, req.Signal), nil
	case CheckTests, "":
		if a.testsDriver == nil {
			return Verdict{Err: a.lang + " function-call grading isn't wired yet — these exercises use stdout/compile checks"}, nil
		}
		driver, err := a.testsDriver(req)
		if err != nil {
			return Verdict{}, err
		}
		if err := r.write(a.file, driver); err != nil {
			return Verdict{}, err
		}
		last, _ := execSteps(r, a.runCmds(r))
		if last.timedOut {
			return Verdict{Err: "time limit exceeded"}, nil
		}
		if a.inLang {
			return ParseInLangOutput(last.stdout+"\n"+last.stderr, req.Tests), nil
		}
		return ParseHarnessOutput(last.stdout+"\n"+last.stderr, req.Tests), nil
	default:
		return Verdict{Err: "unsupported check for " + a.lang + ": " + string(req.Check)}, nil
	}
}

// execSteps runs commands in order. In RUN mode a non-final (compile) step that
// exits nonzero is a compile error, returned as the second value; the LAST step's
// output is returned for the caller to interpret.
func execSteps(r runner, cmds [][]string) (execOut, string) {
	var last execOut
	for i, c := range cmds {
		if len(c) == 0 {
			continue
		}
		last = r.run(c[0], c[1:]...)
		if last.timedOut {
			return last, "time limit exceeded"
		}
		if last.exit != 0 && i < len(cmds)-1 {
			return last, firstLine(last.stderr + last.stdout)
		}
	}
	return last, ""
}

// binName is the platform-correct path of a freshly-built binary in dir.
func binName(dir, base string) string {
	if runtime.GOOS == "windows" {
		base += ".exe"
	}
	return filepath.Join(dir, base)
}

// nativeAdapters returns the per-language oracle-check adapters. C++ is omitted
// until a compiler is commonly available; add it the same way when wired.
func nativeAdapters() map[string]langAdapter {
	return map[string]langAdapter{
		"javascript": nativeAdapter{
			lang: "javascript", file: "main.js",
			runCmds:  func(r runner) [][]string { return [][]string{{r.resolve("node"), "main.js"}} },
			compCmds: func(r runner) [][]string { return [][]string{{r.resolve("node"), "--check", "main.js"}} },
			testsDriver: func(req GradeRequest) (string, error) {
				return BuildJSDriver(req.Source, req.FuncName, req.Tests, req.Shape)
			},
		},
		"typescript": nativeAdapter{
			lang: "typescript", file: "main.ts",
			// --target es2017 so modern built-ins (Set, Map, Object.entries, …) are in
			// lib; bare tsc defaults to ES5, where `new Set()` is "Cannot find name 'Set'".
			// --strict false PINNED: TypeScript 6.0 made strict the default, which
			// rejects the harness driver (untyped __out/__run) and loosely-typed
			// player code on every grade. The function-call grading contract is
			// non-strict; strictness stays explicit in the compiles check below.
			runCmds: func(r runner) [][]string {
				return [][]string{{r.resolve("tsc"), "--strict", "false", "--target", "es2017", "main.ts"}, {r.resolve("node"), "main.js"}}
			},
			testsDriver: func(req GradeRequest) (string, error) {
				return BuildTSDriver(req.Source, req.FuncName, req.Tests, req.Shape)
			},
			// --strict so type-safety exercises (narrowing, any/unknown, exhaustiveness)
			// actually fail as authored — most TS type errors only fire under strict.
			compCmds: func(r runner) [][]string {
				return [][]string{{r.resolve("tsc"), "--noEmit", "--strict", "--target", "es2017", "main.ts"}}
			},
		},
		"go": nativeAdapter{
			lang: "go", file: "main.go",
			runCmds: func(r runner) [][]string { return [][]string{{r.resolve("go"), "run", "main.go"}} },
			compCmds: func(r runner) [][]string {
				return [][]string{{r.resolve("go"), "build", "-o", binName(r.dir, "out"), "main.go"}}
			},
			testsDriver: func(req GradeRequest) (string, error) {
				return BuildGoDriver(req.Source, req.FuncName, req.Tests, req.Shape)
			},
		},
		"java": nativeAdapter{
			lang: "java", file: "Main.java",
			runCmds: func(r runner) [][]string {
				return [][]string{{r.resolve("javac"), "Main.java"}, {r.resolve("java"), "Main"}}
			},
			compCmds: func(r runner) [][]string { return [][]string{{r.resolve("javac"), "Main.java"}} },
			testsDriver: func(req GradeRequest) (string, error) {
				return BuildJavaDriver(req.Source, req.FuncName, req.Tests, req.Shape)
			},
			inLang: true,
		},
		"csharp": nativeAdapter{
			lang: "csharp", file: "Program.cs",
			extra: func(dir string) map[string]string {
				return map[string]string{"app.csproj": csprojTemplate}
			},
			runCmds:  func(r runner) [][]string { return [][]string{{r.resolve("dotnet"), "run", "--nologo"}} },
			compCmds: func(r runner) [][]string { return [][]string{{r.resolve("dotnet"), "build", "--nologo"}} },
			testsDriver: func(req GradeRequest) (string, error) {
				return BuildCSharpDriver(req.Source, req.FuncName, req.Tests, req.Shape)
			},
		},
		"rust": nativeAdapter{
			lang: "rust", file: "main.rs",
			runCmds: func(r runner) [][]string {
				out := binName(r.dir, "out")
				return [][]string{{r.resolve("rustc"), "--edition", "2021", "main.rs", "-o", out}, {out}}
			},
			compCmds: func(r runner) [][]string {
				return [][]string{{r.resolve("rustc"), "--edition", "2021", "main.rs", "-o", binName(r.dir, "out")}}
			},
			testsDriver: func(req GradeRequest) (string, error) {
				return BuildRustDriver(req.Source, req.FuncName, req.Tests, req.Shape)
			},
			inLang: true,
		},
	}
}

// csprojTemplate is a minimal console project for the C# adapter. TargetFramework
// is pinned to a known-good language baseline; RollForward lets the output RUN
// on whatever newer runtime the player actually has (BYO toolchain — players
// with only .NET 10/11 must not fail on a net9.0 target; the TS 6.0 lesson).
const csprojTemplate = `<Project Sdk="Microsoft.NET.Sdk">
  <PropertyGroup>
    <OutputType>Exe</OutputType>
    <TargetFramework>net9.0</TargetFramework>
    <RollForward>LatestMajor</RollForward>
    <ImplicitUsings>enable</ImplicitUsings>
    <Nullable>disable</Nullable>
    <AssemblyName>devascentapp</AssemblyName>
  </PropertyGroup>
</Project>
`
