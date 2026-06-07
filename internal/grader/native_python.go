package grader

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// NativePython runs player Python via the local interpreter.
//
// ⚠️ DEV-ONLY: this backend is NOT sandboxed (it runs arbitrary Python on the
// host) and therefore violates NFR-1. WazeroPython is now the default (grader.New);
// this backend is the opt-in dev-speed escape hatch (DEVASCENT_GRADER=native) and
// the reference oracle for the wazero equivalence gate. Do not ship as default.
type NativePython struct {
	Exe     string
	Timeout time.Duration
}

func NewNativePython() NativePython {
	exe := os.Getenv("DEVASCENT_PYTHON")
	if exe == "" {
		exe = "python"
	}
	return NativePython{Exe: exe, Timeout: 10 * time.Second}
}

const marker = "__DEVASCENT_RESULT__"

func (g NativePython) Run(lang, source, funcName string, tests []TestCase, shape Shape) (Verdict, error) {
	if lang != "python" {
		return Verdict{Err: "unsupported language: " + lang}, nil
	}
	driver, err := BuildPyDriver(source, funcName, tests, shape)
	if err != nil {
		return Verdict{}, err
	}

	dir, err := os.MkdirTemp("", "devascent-grade-")
	if err != nil {
		return Verdict{}, err
	}
	defer os.RemoveAll(dir)
	file := filepath.Join(dir, "run.py")
	if err := os.WriteFile(file, []byte(driver), 0o600); err != nil {
		return Verdict{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), g.Timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, g.Exe, file)
	cmd.Dir = dir
	out, _ := cmd.CombinedOutput()

	if ctx.Err() == context.DeadlineExceeded {
		return Verdict{Err: "time limit exceeded"}, nil
	}
	return ParseHarnessOutput(string(out), tests), nil
}

// Grade satisfies the Grader interface. NativePython is the dev-only Python
// oracle and supports only the tests check; the compile-error/compiles/stdout
// modes are the LocalToolchain backend's job.
func (g NativePython) Grade(req GradeRequest) (Verdict, error) {
	switch req.Check {
	case CheckTests, "":
		return g.Run(req.Lang, req.Source, req.FuncName, req.Tests, req.Shape)
	default:
		return Verdict{Err: "NativePython supports only the tests check, got: " + string(req.Check)}, nil
	}
}

// BuildPyDriver assembles the full Python program to execute: shape prelude +
// player source + the JSON-in/JSON-out harness. Backend-agnostic — both the
// native and (future) wazero backends run THIS exact string, so grading stays
// identical regardless of how Python executes.
func BuildPyDriver(source, funcName string, tests []TestCase, shape Shape) (string, error) {
	testsJSON, err := json.Marshal(tests)
	if err != nil {
		return "", err
	}
	b64 := base64.StdEncoding.EncodeToString(testsJSON)
	return shapePrelude(shape.Kind) + source + "\n\n" + pyHarness(funcName, b64, shape), nil
}

// ParseHarnessOutput turns the program's stdout into a Verdict (find the marker
// line, compare each case with float tolerance). Backend-agnostic.
func ParseHarnessOutput(out string, tests []TestCase) Verdict {
	line := extractMarker(out)
	if line == "" {
		return Verdict{Err: trimErr(out)}
	}
	var raw []struct {
		Name  string          `json:"name"`
		Got   json.RawMessage `json:"got"`
		Error string          `json:"error"`
	}
	if err := json.Unmarshal([]byte(line), &raw); err != nil {
		return Verdict{Err: "bad harness output: " + err.Error()}
	}
	byName := map[string]TestCase{}
	for _, t := range tests {
		byName[t.Name] = t
	}
	v := Verdict{Passed: true}
	for _, r := range raw {
		tc := byName[r.Name]
		exp, _ := json.Marshal(tc.Expected)
		cr := CaseResult{Name: r.Name, Expected: string(exp), Got: string(r.Got)}
		switch {
		case r.Error != "":
			cr.Passed, cr.Err = false, r.Error
		case equalJSON(r.Got, tc.Expected):
			cr.Passed = true
		default:
			cr.Passed = false
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

func pyHarness(funcName, testsB64 string, shape Shape) string {
	argKindsPy := "[]"
	if len(shape.ArgKinds) > 0 {
		parts := make([]string, len(shape.ArgKinds))
		for i, k := range shape.ArgKinds {
			parts[i] = strconv.Quote(k)
		}
		argKindsPy = "[" + strings.Join(parts, ", ") + "]"
	}
	retNode := "False"
	if shape.RetKind == "node" {
		retNode = "True"
	}
	return fmt.Sprintf(`import json as __json, base64 as __b64
__tests = __json.loads(__b64.b64decode(%q).decode("utf-8"))
__argkinds = %s
__retnode = %s
__out = []
for __t in __tests:
    try:
        __args = []
        for __i, __a in enumerate(__t["input"]):
            if __i < len(__argkinds) and __argkinds[__i] == "node":
                __args.append(__build(__a))
            else:
                __args.append(__a)
        __r = %s(*__args)
        if __retnode:
            __r = __dump(__r)
        __out.append({"name": __t["name"], "got": __r})
    except Exception as __e:
        __out.append({"name": __t["name"], "error": str(__e)})
print(%q + __json.dumps(__out))
`, testsB64, argKindsPy, retNode, funcName, marker)
}

// shapePrelude returns the Python class + __build/__dump helpers for a node
// shape, injected before the player source. NFR-1: this prelude is the
// native-Python backend's; the wazero/WASM backend must reimplement it to keep
// data-structure (linked-list/tree) problem support.
func shapePrelude(kind string) string {
	switch kind {
	case "linkedlist":
		return `class ListNode:
    def __init__(self, val=0, next=None):
        self.val = val
        self.next = next
def __build(arr):
    head = None
    for __v in reversed(arr):
        head = ListNode(__v, head)
    return head
def __dump(node):
    out = []
    while node is not None:
        out.append(node.val)
        node = node.next
    return out
`
	case "tree":
		return `from collections import deque as __deque
class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right
def __build(arr):
    if not arr:
        return None
    __it = iter(arr)
    root = TreeNode(next(__it))
    __q = __deque([root])
    while __q:
        __node = __q.popleft()
        try:
            __lv = next(__it)
        except StopIteration:
            break
        if __lv is not None:
            __node.left = TreeNode(__lv); __q.append(__node.left)
        try:
            __rv = next(__it)
        except StopIteration:
            break
        if __rv is not None:
            __node.right = TreeNode(__rv); __q.append(__node.right)
    return root
def __dump(root):
    if root is None:
        return []
    out = []
    __q = __deque([root])
    while __q:
        __node = __q.popleft()
        if __node is None:
            out.append(None)
        else:
            out.append(__node.val)
            __q.append(__node.left); __q.append(__node.right)
    while out and out[-1] is None:
        out.pop()
    return out
`
	}
	return ""
}

func extractMarker(out string) string {
	for _, ln := range strings.Split(out, "\n") {
		if i := strings.Index(ln, marker); i >= 0 {
			return strings.TrimSpace(ln[i+len(marker):])
		}
	}
	return ""
}

func equalJSON(got json.RawMessage, expected any) bool {
	eb, err := json.Marshal(expected)
	if err != nil {
		return false
	}
	var a, b any
	if err := json.Unmarshal(got, &a); err != nil {
		return false
	}
	if err := json.Unmarshal(eb, &b); err != nil {
		return false
	}
	return jsonEqual(a, b)
}

// jsonEqual compares two decoded-JSON values with NUMERIC TOLERANCE for floats
// (so 1.0000000001 == 1.0 and tax/median rounding doesn't cause false fails),
// recursing through arrays and objects. Ints (which JSON-decode to float64),
// strings, bools, and null stay effectively exact.
func jsonEqual(a, b any) bool {
	switch av := a.(type) {
	case float64:
		bv, ok := b.(float64)
		return ok && floatsClose(av, bv)
	case []any:
		bv, ok := b.([]any)
		if !ok || len(av) != len(bv) {
			return false
		}
		for i := range av {
			if !jsonEqual(av[i], bv[i]) {
				return false
			}
		}
		return true
	case map[string]any:
		bv, ok := b.(map[string]any)
		if !ok || len(av) != len(bv) {
			return false
		}
		for k, va := range av {
			vb, ok := bv[k]
			if !ok || !jsonEqual(va, vb) {
				return false
			}
		}
		return true
	default:
		return reflect.DeepEqual(a, b)
	}
}

func floatsClose(a, b float64) bool {
	diff := math.Abs(a - b)
	if diff <= 1e-9 {
		return true
	}
	scale := math.Max(1, math.Max(math.Abs(a), math.Abs(b)))
	return diff <= 1e-6*scale
}

func trimErr(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "no output from program"
	}
	return s
}
