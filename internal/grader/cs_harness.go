package grader

import (
	"fmt"
	"strconv"
	"strings"
)

// cs_harness.go — Model A function-call grading for C#. Mirrors the Go harness
// (stdlib System.Text.Json for serialization), reusing the shared gtype inference
// (inferType/argTypesFor). The player writes a `public class Solution { public
// RetType FuncName(...) {...} }`; the harness adds a Program that instantiates it,
// calls it per test with embedded literals, serializes the result, and prints the
// shared marker line for ParseHarnessOutput + jsonEqual.

func csTypeStr(t gtype) string {
	switch t.kind {
	case "int":
		return "long"
	case "float":
		return "double"
	case "string":
		return "string"
	case "bool":
		return "bool"
	case "slice":
		return csTypeStr(deref(t.elem)) + "[]"
	case "list":
		return "ListNode"
	case "tree":
		return "TreeNode"
	default:
		return "object"
	}
}

func csLiteral(t gtype, v any) string {
	switch t.kind {
	case "int":
		return goScalar(v) // bare integer literal, implicitly convertible to long
	case "float":
		if f, ok := asFloat(v); ok {
			return strconv.FormatFloat(f, 'g', -1, 64) + "d"
		}
		return "0d"
	case "string":
		s, _ := v.(string)
		return strconv.Quote(s)
	case "bool":
		b, _ := v.(bool)
		return strconv.FormatBool(b)
	case "slice":
		arr, _ := v.([]any)
		parts := make([]string, len(arr))
		for i, e := range arr {
			parts[i] = csLiteral(deref(t.elem), e)
		}
		return "new " + csTypeStr(t) + "{" + strings.Join(parts, ", ") + "}"
	case "any":
		return csAny(v)
	default:
		return "null"
	}
}

// csAny renders a heterogeneous (object) input value as a C# literal — the
// analogue of goAny. Without this, op-list/design inputs (e.g. [["push",-2],
// ["pop"]]) collapsed to null. Integers box as long (5L) so the player's
// (long)op[i] casts succeed. (Output comparison is shared-JSON, so only the
// INPUT side needs this.)
func csAny(v any) string {
	switch x := v.(type) {
	case nil:
		return "null"
	case bool:
		return strconv.FormatBool(x)
	case int, int64, int32:
		return goScalar(x) + "L"
	case float64:
		if x == float64(int64(x)) {
			return strconv.FormatInt(int64(x), 10) + "L"
		}
		return strconv.FormatFloat(x, 'g', -1, 64) + "d"
	case string:
		return strconv.Quote(x)
	case []any:
		parts := make([]string, len(x))
		for i, e := range x {
			parts[i] = csAny(e)
		}
		return "new object[]{" + strings.Join(parts, ", ") + "}"
	default:
		return "null"
	}
}

// csNodeBuildLit emits BuildList/BuildTree(<array literal>) for a node argument.
func csNodeBuildLit(shapeKind string, v any) string {
	arr, _ := v.([]any)
	if shapeKind == "tree" {
		parts := make([]string, len(arr))
		for i, e := range arr {
			if e == nil {
				parts[i] = "null"
			} else {
				parts[i] = goScalar(e)
			}
		}
		return "BuildTree(new int?[]{" + strings.Join(parts, ", ") + "})"
	}
	parts := make([]string, len(arr))
	for i, e := range arr {
		parts[i] = goScalar(e)
	}
	return "BuildList(new int[]{" + strings.Join(parts, ", ") + "})"
}

func csNodeClasses(shapeKind string) string {
	if shapeKind == "linkedlist" {
		return "public class ListNode { public int val; public ListNode next; public ListNode(int v){ val = v; next = null; } }\n\n"
	}
	if shapeKind == "tree" {
		return "public class TreeNode { public int val; public TreeNode left; public TreeNode right; public TreeNode(int v){ val = v; } }\n\n"
	}
	return ""
}

func csNodeHelpers(shape Shape) string {
	var b strings.Builder
	if shape.Kind == "linkedlist" {
		b.WriteString("    static ListNode BuildList(int[] a){ ListNode head = null; for (int i = a.Length - 1; i >= 0; i--){ var n = new ListNode(a[i]); n.next = head; head = n; } return head; }\n")
		if shape.RetKind == "node" {
			b.WriteString("    static long[] DumpList(ListNode n){ var o = new List<long>(); while (n != null){ o.Add(n.val); n = n.next; } return o.ToArray(); }\n")
		}
	}
	if shape.Kind == "tree" {
		b.WriteString("    static TreeNode BuildTree(int?[] a){ if (a.Length == 0 || a[0] == null) return null; var root = new TreeNode(a[0].Value); var q = new Queue<TreeNode>(); q.Enqueue(root); int i = 1; while (q.Count > 0 && i < a.Length){ var node = q.Dequeue(); if (i < a.Length){ if (a[i] != null){ node.left = new TreeNode(a[i].Value); q.Enqueue(node.left); } i++; } if (i < a.Length){ if (a[i] != null){ node.right = new TreeNode(a[i].Value); q.Enqueue(node.right); } i++; } } return root; }\n")
		if shape.RetKind == "node" {
			b.WriteString("    static object DumpTree(TreeNode root){ var o = new List<object>(); if (root == null) return o; var q = new Queue<TreeNode>(); q.Enqueue(root); while (q.Count > 0){ var n = q.Dequeue(); if (n == null){ o.Add(null); continue; } o.Add((long)n.val); q.Enqueue(n.left); q.Enqueue(n.right); } while (o.Count > 0 && o[o.Count - 1] == null) o.RemoveAt(o.Count - 1); return o; }\n")
		}
	}
	return b.String()
}

// BuildCSharpDriver assembles the full C# program.
func BuildCSharpDriver(source, funcName string, tests []TestCase, shape Shape) (string, error) {
	argTypes := argTypesFor(tests, shape)
	var b strings.Builder
	b.WriteString("using System;\nusing System.Collections.Generic;\nusing System.Text.Json;\n\n")
	if shape.Kind != "" {
		b.WriteString(csNodeClasses(shape.Kind))
	}
	b.WriteString(strings.TrimRight(source, "\n") + "\n\n")
	b.WriteString("class __Harness {\n")
	b.WriteString(csNodeHelpers(shape))
	b.WriteString("    static void __Run(List<object> outl, string name, Func<object> f){ try { outl.Add(new { name = name, got = f() }); } catch (Exception e){ outl.Add(new { name = name, error = e.Message }); } }\n")
	b.WriteString("    static void Main(){\n")
	b.WriteString("        var __out = new List<object>();\n")
	for _, tc := range tests {
		argStrs := make([]string, len(argTypes))
		for i := range argTypes {
			var v any
			if i < len(tc.Input) {
				v = tc.Input[i]
			}
			if i < len(shape.ArgKinds) && shape.ArgKinds[i] == "node" {
				argStrs[i] = csNodeBuildLit(shape.Kind, v)
			} else {
				argStrs[i] = csLiteral(argTypes[i], v)
			}
		}
		call := "new Solution()." + funcName + "(" + strings.Join(argStrs, ", ") + ")"
		if shape.RetKind == "node" {
			if shape.Kind == "tree" {
				call = "DumpTree(" + call + ")"
			} else {
				call = "DumpList(" + call + ")"
			}
		}
		b.WriteString("        __Run(__out, " + strconv.Quote(tc.Name) + ", () => " + call + ");\n")
	}
	b.WriteString("        Console.WriteLine(" + strconv.Quote(marker) + " + JsonSerializer.Serialize(__out));\n")
	b.WriteString("    }\n}\n")
	return b.String(), nil
}

func csZero(t gtype) string {
	// `default` works for any C# type (0, false, null), so the stub always compiles.
	return "default"
}

// CSharpStarter generates a Solution class stub for a problem.
func CSharpStarter(funcName, pySource string, tests []TestCase, shape Shape) string {
	argTypes := argTypesFor(tests, shape)
	names := paramNames(pySource, len(argTypes))
	params := make([]string, len(argTypes))
	for i, t := range argTypes {
		params[i] = csTypeStr(t) + " " + names[i]
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
	return fmt.Sprintf("public class Solution {\n    public %s %s(%s) {\n        // your code here\n        return %s;\n    }\n}\n",
		csTypeStr(ret), funcName, strings.Join(params, ", "), csZero(ret))
}
