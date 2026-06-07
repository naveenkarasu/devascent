package grader

import (
	"fmt"
	"strconv"
	"strings"
)

// java_harness.go — Model A function-call grading for Java via IN-LANGUAGE
// comparison (no stdlib JSON). The player writes `class Solution { public Ret
// funcName(...) {...} }`; the harness (public class Main) embeds args + the
// expected value as Java literals, compares with ==/equals/Arrays.equals, and
// prints the in-language line protocol (ParseInLangOutput). Node returns are
// dumped to arrays and compared.

func javaTypeStr(t gtype) string {
	switch t.kind {
	case "int":
		return "long"
	case "float":
		return "double"
	case "string":
		return "String"
	case "bool":
		return "boolean"
	case "slice":
		return javaTypeStr(deref(t.elem)) + "[]"
	case "list":
		return "ListNode"
	case "tree":
		return "TreeNode"
	default:
		return "Object"
	}
}

func javaLiteral(t gtype, v any) string {
	switch t.kind {
	case "int":
		return goScalar(v) + "L"
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
			parts[i] = javaLiteral(deref(t.elem), e)
		}
		return "new " + javaTypeStr(t) + "{" + strings.Join(parts, ", ") + "}"
	default:
		return "null"
	}
}

func javaNodeBuildLit(shapeKind string, v any) string {
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
		return "buildTree(new Integer[]{" + strings.Join(parts, ", ") + "})"
	}
	parts := make([]string, len(arr))
	for i, e := range arr {
		parts[i] = goScalar(e)
	}
	return "buildList(new int[]{" + strings.Join(parts, ", ") + "})"
}

func javaNodeClasses(shapeKind string) string {
	if shapeKind == "linkedlist" {
		return "class ListNode { int val; ListNode next; ListNode(int v){ val = v; } }\n\n"
	}
	if shapeKind == "tree" {
		return "class TreeNode { int val; TreeNode left, right; TreeNode(int v){ val = v; } }\n\n"
	}
	return ""
}

func javaNodeHelpers(shape Shape) string {
	var b strings.Builder
	if shape.Kind == "linkedlist" {
		b.WriteString("  static ListNode buildList(int[] a){ ListNode head = null; for (int i = a.length - 1; i >= 0; i--){ ListNode n = new ListNode(a[i]); n.next = head; head = n; } return head; }\n")
		if shape.RetKind == "node" {
			b.WriteString("  static long[] dumpList(ListNode n){ java.util.List<Long> o = new java.util.ArrayList<>(); while (n != null){ o.add((long)n.val); n = n.next; } long[] r = new long[o.size()]; for (int i = 0; i < r.length; i++) r[i] = o.get(i); return r; }\n")
		}
	}
	if shape.Kind == "tree" {
		b.WriteString("  static TreeNode buildTree(Integer[] a){ if (a.length == 0 || a[0] == null) return null; TreeNode root = new TreeNode(a[0]); java.util.Queue<TreeNode> q = new java.util.LinkedList<>(); q.add(root); int i = 1; while (!q.isEmpty() && i < a.length){ TreeNode node = q.poll(); if (i < a.length){ if (a[i] != null){ node.left = new TreeNode(a[i]); q.add(node.left); } i++; } if (i < a.length){ if (a[i] != null){ node.right = new TreeNode(a[i]); q.add(node.right); } i++; } } return root; }\n")
		if shape.RetKind == "node" {
			b.WriteString("  static Long[] dumpTree(TreeNode root){ java.util.List<Long> o = new java.util.ArrayList<>(); if (root == null) return new Long[0]; java.util.Queue<TreeNode> q = new java.util.LinkedList<>(); q.add(root); while (!q.isEmpty()){ TreeNode n = q.poll(); if (n == null){ o.add(null); continue; } o.add((long)n.val); q.add(n.left); q.add(n.right); } while (!o.isEmpty() && o.get(o.size() - 1) == null) o.remove(o.size() - 1); return o.toArray(new Long[0]); }\n")
		}
	}
	return b.String()
}

func javaRetInfo(tests []TestCase, shape Shape) (ret gtype, isList, isTree bool) {
	if shape.RetKind == "node" {
		if shape.Kind == "tree" {
			return gtype{kind: "any"}, false, true
		}
		return gtype{kind: "any"}, true, false
	}
	var samples []any
	for _, tc := range tests {
		samples = append(samples, tc.Expected)
	}
	return inferType(samples), false, false
}

func javaExpectedLit(ret gtype, v any, isList, isTree bool) string {
	if isList {
		arr, _ := v.([]any)
		parts := make([]string, len(arr))
		for i, e := range arr {
			parts[i] = goScalar(e) + "L"
		}
		return "new long[]{" + strings.Join(parts, ", ") + "}"
	}
	if isTree {
		arr, _ := v.([]any)
		parts := make([]string, len(arr))
		for i, e := range arr {
			if e == nil {
				parts[i] = "null"
			} else {
				parts[i] = goScalar(e) + "L"
			}
		}
		return "new Long[]{" + strings.Join(parts, ", ") + "}"
	}
	return javaLiteral(ret, v)
}

func javaCompare(ret gtype, got, exp string, isList, isTree bool) string {
	if isList || isTree {
		return "java.util.Arrays.equals(" + got + ", " + exp + ")"
	}
	switch ret.kind {
	case "int", "bool":
		return got + " == " + exp
	case "float":
		return "Math.abs(" + got + " - " + exp + ") < 1e-6"
	case "string":
		return exp + ".equals(" + got + ")"
	case "slice":
		if deref(ret.elem).kind == "slice" {
			return "java.util.Arrays.deepEquals(" + got + ", " + exp + ")"
		}
		return "java.util.Arrays.equals(" + got + ", " + exp + ")"
	default:
		return "String.valueOf(" + got + ").equals(String.valueOf(" + exp + "))"
	}
}

func javaRepr(got string, ret gtype, isList, isTree bool) string {
	if isList || isTree {
		return "java.util.Arrays.toString(" + got + ")"
	}
	if ret.kind == "slice" {
		if deref(ret.elem).kind == "slice" {
			return "java.util.Arrays.deepToString(" + got + ")"
		}
		return "java.util.Arrays.toString(" + got + ")"
	}
	return "String.valueOf(" + got + ")"
}

// BuildJavaDriver assembles the full Java program (node classes + player Solution
// + public class Main harness).
func BuildJavaDriver(source, funcName string, tests []TestCase, shape Shape) (string, error) {
	argTypes := argTypesFor(tests, shape)
	ret, isList, isTree := javaRetInfo(tests, shape)
	var b strings.Builder
	if shape.Kind != "" {
		b.WriteString(javaNodeClasses(shape.Kind))
	}
	b.WriteString(strings.TrimRight(source, "\n") + "\n\n")
	b.WriteString("public class Main {\n")
	b.WriteString(javaNodeHelpers(shape))
	b.WriteString("  public static void main(String[] __a){\n")
	b.WriteString("    final String M = " + strconv.Quote(marker) + "; final String US = \"\\u001f\";\n")
	for _, tc := range tests {
		args := make([]string, len(argTypes))
		for i := range argTypes {
			var v any
			if i < len(tc.Input) {
				v = tc.Input[i]
			}
			if i < len(shape.ArgKinds) && shape.ArgKinds[i] == "node" {
				args[i] = javaNodeBuildLit(shape.Kind, v)
			} else {
				args[i] = javaLiteral(argTypes[i], v)
			}
		}
		call := "new Solution()." + funcName + "(" + strings.Join(args, ", ") + ")"
		gotType, gotExpr := javaTypeStr(ret), call
		if isList {
			gotType, gotExpr = "long[]", "dumpList("+call+")"
		} else if isTree {
			gotType, gotExpr = "Long[]", "dumpTree("+call+")"
		}
		exp := javaExpectedLit(ret, tc.Expected, isList, isTree)
		cmp := javaCompare(ret, "__got", exp, isList, isTree)
		repr := javaRepr("__got", ret, isList, isTree)
		q := strconv.Quote(tc.Name)
		b.WriteString("    try {\n")
		b.WriteString("      " + gotType + " __got = " + gotExpr + ";\n")
		b.WriteString("      boolean __ok = " + cmp + ";\n")
		b.WriteString("      System.out.println(M + " + q + " + US + (__ok ? \"1\" : \"0\") + US + " + repr + ");\n")
		b.WriteString("    } catch (Throwable __e) {\n")
		b.WriteString("      System.out.println(M + " + q + " + US + \"0\" + US + (\"ex: \" + __e));\n")
		b.WriteString("    }\n")
	}
	b.WriteString("  }\n}\n")
	return b.String(), nil
}

// JavaStarter generates a Solution class stub for a problem.
func JavaStarter(funcName, pySource string, tests []TestCase, shape Shape) string {
	argTypes := argTypesFor(tests, shape)
	names := paramNames(pySource, len(argTypes))
	params := make([]string, len(argTypes))
	for i, t := range argTypes {
		params[i] = javaTypeStr(t) + " " + names[i]
	}
	var ret gtype
	if shape.RetKind == "node" {
		ret = nodeType(shape.Kind)
	} else {
		ret, _, _ = javaRetInfo(tests, shape)
	}
	return fmt.Sprintf("class Solution {\n    public %s %s(%s) {\n        // your code here\n        return %s;\n    }\n}\n",
		javaTypeStr(ret), funcName, strings.Join(params, ", "), javaZero(ret))
}

func javaZero(t gtype) string {
	switch t.kind {
	case "int":
		return "0L"
	case "float":
		return "0d"
	case "bool":
		return "false"
	default:
		return "null"
	}
}
