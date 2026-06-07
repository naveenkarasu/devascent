package grader

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// js_harness.go — Model A function-call grading for JavaScript and TypeScript.
// JS is dynamically typed, so args are embedded as JSON literals (JSON ⊂ JS) and
// the result is JSON.stringify'd — reusing ParseHarnessOutput + jsonEqual exactly
// like Go/C#. No type inference is needed for JS. TypeScript shares the harness
// (the call is cast `as any`) but compiles with tsc first; its starter stub IS
// typed (inferred) for the nicer authoring experience.

// jsLit renders any JSON value as a JS literal (JSON is valid JS for our data).
func jsLit(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return "null"
	}
	return string(b)
}

func jsNodeBuildLit(shapeKind string, v any) string {
	arr := jsLit(v) // a JS array literal (with nulls for tree)
	if shapeKind == "tree" {
		return "buildTree(" + arr + ")"
	}
	return "buildList(" + arr + ")"
}

func jsPrelude(shape Shape) string {
	var b strings.Builder
	if shape.Kind == "linkedlist" {
		b.WriteString("class ListNode { val; next; constructor(val, next = null) { this.val = val; this.next = next; } }\n")
		b.WriteString("function buildList(a) { let head = null; for (let i = a.length - 1; i >= 0; i--) head = new ListNode(a[i], head); return head; }\n")
		if shape.RetKind == "node" {
			b.WriteString("function dumpList(n) { const o = []; while (n) { o.push(n.val); n = n.next; } return o; }\n")
		}
	}
	if shape.Kind == "tree" {
		b.WriteString("class TreeNode { val; left; right; constructor(val, left = null, right = null) { this.val = val; this.left = left; this.right = right; } }\n")
		b.WriteString("function buildTree(a) { if (!a.length || a[0] === null) return null; const root = new TreeNode(a[0]); const q = [root]; let i = 1; while (q.length && i < a.length) { const node = q.shift(); if (i < a.length) { if (a[i] !== null) { node.left = new TreeNode(a[i]); q.push(node.left); } i++; } if (i < a.length) { if (a[i] !== null) { node.right = new TreeNode(a[i]); q.push(node.right); } i++; } } return root; }\n")
		if shape.RetKind == "node" {
			b.WriteString("function dumpTree(root) { const o = []; if (!root) return o; const q = [root]; while (q.length) { const n = q.shift(); if (n === null) { o.push(null); continue; } o.push(n.val); q.push(n.left, n.right); } while (o.length && o[o.length - 1] === null) o.pop(); return o; }\n")
		}
	}
	if b.Len() > 0 {
		b.WriteString("\n")
	}
	return b.String()
}

// buildJSDriver assembles the JS/TS program. asAny casts the call to bypass TS
// type-checking on the harness side (the player's function still type-checks).
func buildJSDriver(source, funcName string, tests []TestCase, shape Shape, asAny bool) string {
	var b strings.Builder
	b.WriteString(jsPrelude(shape))
	b.WriteString(strings.TrimRight(source, "\n") + "\n\n")
	b.WriteString("const __out = [];\n")
	b.WriteString("function __run(name, f) { try { __out.push({ name: name, got: f() }); } catch (e) { __out.push({ name: name, error: String((e && e.message) || e) }); } }\n")
	callee := funcName
	if asAny {
		callee = "(" + funcName + " as any)"
	}
	nargs := 0
	if len(tests) > 0 {
		nargs = len(tests[0].Input)
	}
	for _, tc := range tests {
		args := make([]string, nargs)
		for i := 0; i < nargs; i++ {
			var v any
			if i < len(tc.Input) {
				v = tc.Input[i]
			}
			if i < len(shape.ArgKinds) && shape.ArgKinds[i] == "node" {
				args[i] = jsNodeBuildLit(shape.Kind, v)
			} else {
				args[i] = jsLit(v)
			}
		}
		call := callee + "(" + strings.Join(args, ", ") + ")"
		if shape.RetKind == "node" {
			if shape.Kind == "tree" {
				call = "dumpTree(" + call + ")"
			} else {
				call = "dumpList(" + call + ")"
			}
		}
		b.WriteString("__run(" + strconv.Quote(tc.Name) + ", () => " + call + ");\n")
	}
	b.WriteString("console.log(" + strconv.Quote(marker) + " + JSON.stringify(__out));\n")
	return b.String()
}

// BuildJSDriver builds the JavaScript grading program.
func BuildJSDriver(source, funcName string, tests []TestCase, shape Shape) (string, error) {
	return buildJSDriver(source, funcName, tests, shape, false), nil
}

// BuildTSDriver builds the TypeScript grading program (call cast `as any`).
func BuildTSDriver(source, funcName string, tests []TestCase, shape Shape) (string, error) {
	return buildJSDriver(source, funcName, tests, shape, true), nil
}

func jsArgNames(pySource string, n int) []string { return paramNames(pySource, n) }

// JSStarter generates a JavaScript function stub (untyped).
func JSStarter(funcName, pySource string, tests []TestCase, shape Shape) string {
	n := 0
	if len(tests) > 0 {
		n = len(tests[0].Input)
	}
	names := jsArgNames(pySource, n)
	return fmt.Sprintf("function %s(%s) {\n  // your code here\n}\n", funcName, strings.Join(names, ", "))
}

func tsTypeStr(t gtype) string {
	switch t.kind {
	case "int", "float":
		return "number"
	case "string":
		return "string"
	case "bool":
		return "boolean"
	case "slice":
		return tsTypeStr(deref(t.elem)) + "[]"
	case "list":
		return "ListNode"
	case "tree":
		return "TreeNode"
	default:
		return "any"
	}
}

// TSStarter generates a typed TypeScript function stub (inferred types).
func TSStarter(funcName, pySource string, tests []TestCase, shape Shape) string {
	argTypes := argTypesFor(tests, shape)
	names := paramNames(pySource, len(argTypes))
	params := make([]string, len(argTypes))
	for i, t := range argTypes {
		params[i] = names[i] + ": " + tsTypeStr(t)
	}
	ret := gtype{kind: "any"}
	if shape.RetKind == "node" {
		ret = nodeType(shape.Kind)
	} else {
		var samples []any
		for _, tc := range tests {
			samples = append(samples, tc.Expected)
		}
		ret = inferType(samples)
	}
	return fmt.Sprintf("function %s(%s): %s {\n  // your code here\n  return undefined as any\n}\n",
		funcName, strings.Join(params, ", "), tsTypeStr(ret))
}
