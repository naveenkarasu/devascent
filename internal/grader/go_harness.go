package grader

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

// go_harness.go — Model A (LeetCode-style) function-call grading for Go.
//
// Generates a complete Go program: the player's source (defining funcName) + a
// main that embeds each test's args as native Go literals, calls the function,
// json.Marshal's the result, and prints one marker line. The Go-side
// ParseHarnessOutput + jsonEqual then compare to Expected exactly as for Python
// (so float tolerance & canonical-order semantics match). Types are inferred from
// the test data (self-consistent: the same inference drives the starter stub, so
// the player's signature always matches the harness call).

// gtype is an inferred Go type for a value position.
type gtype struct {
	kind string // int | float | string | bool | slice | list | tree | any
	elem *gtype // slice element
}

func nodeType(shapeKind string) gtype {
	if shapeKind == "tree" {
		return gtype{kind: "tree"}
	}
	return gtype{kind: "list"}
}

// inferType merges the inferred type across several sample values (resolving
// empty arrays by looking at other cases).
// inferType merges the inferred type of every sample. The accumulator starts at
// the bottom type "unknown" (no evidence yet); leftover "unknown" (e.g. all
// samples null/empty) resolves to "any" (interface{}).
func inferType(samples []any) gtype {
	cur := gtype{kind: "unknown"}
	for _, s := range samples {
		cur = mergeType(cur, inferOne(s))
	}
	return resolveUnknown(cur)
}

func inferOne(v any) gtype {
	switch x := v.(type) {
	case nil:
		return gtype{kind: "null"} // a null literal; mixing with a concrete type → nullable (interface{})
	case bool:
		return gtype{kind: "bool"}
	case int, int64, int32:
		return gtype{kind: "int"} // YAML decodes whole numbers as int
	case float64:
		if x == math.Trunc(x) && math.Abs(x) < 1e15 {
			return gtype{kind: "int"}
		}
		return gtype{kind: "float"}
	case string:
		return gtype{kind: "string"}
	case []any:
		el := gtype{kind: "unknown"}
		for _, e := range x {
			el = mergeType(el, inferOne(e))
		}
		ec := el
		return gtype{kind: "slice", elem: &ec}
	default:
		return gtype{kind: "any"}
	}
}

// mergeType unifies two inferred types. "unknown" is the bottom (no evidence →
// identity); "any" is the top (genuinely heterogeneous → absorbing, so a mixed
// list like [["push",-2],["pop"]] stays slice<any>=interface{}, not slice<string>).
func mergeType(a, b gtype) gtype {
	if a.kind == "unknown" {
		return b
	}
	if b.kind == "unknown" {
		return a
	}
	if a.kind == "null" && b.kind == "null" {
		return gtype{kind: "null"} // all-null so far; resolved to interface{} at the end
	}
	if a.kind == "null" || b.kind == "null" {
		return gtype{kind: "any"} // concrete + null → nullable, so null is preserved (interface{})
	}
	if a.kind == "any" || b.kind == "any" {
		return gtype{kind: "any"} // top absorbs — mixed stays mixed
	}
	if a.kind == "slice" && b.kind == "slice" {
		m := mergeType(deref(a.elem), deref(b.elem))
		return gtype{kind: "slice", elem: &m}
	}
	if a.kind == b.kind {
		return a
	}
	if (a.kind == "int" && b.kind == "float") || (a.kind == "float" && b.kind == "int") {
		return gtype{kind: "float"}
	}
	return gtype{kind: "any"} // mismatch → interface{}
}

// resolveUnknown maps any leftover bottom "unknown" (and nested unknowns) to
// "any" so a fully-evidence-free type renders as interface{}.
func resolveUnknown(t gtype) gtype {
	switch t.kind {
	case "unknown", "null":
		return gtype{kind: "any"} // no concrete evidence (empty/all-null) → interface{}
	case "slice":
		e := resolveUnknown(deref(t.elem))
		return gtype{kind: "slice", elem: &e}
	}
	return t
}

func deref(t *gtype) gtype {
	if t == nil {
		return gtype{kind: "any"}
	}
	return *t
}

// goTypeStr renders a Go type for a stub signature.
func goTypeStr(t gtype) string {
	switch t.kind {
	case "int":
		return "int"
	case "float":
		return "float64"
	case "string":
		return "string"
	case "bool":
		return "bool"
	case "slice":
		return "[]" + goTypeStr(deref(t.elem))
	case "list":
		return "*ListNode"
	case "tree":
		return "*TreeNode"
	default:
		return "interface{}"
	}
}

// goLiteral renders a Go literal for value v at inferred type t.
func goLiteral(t gtype, v any) string {
	switch t.kind {
	case "int":
		return goScalar(v)
	case "float":
		if f, ok := asFloat(v); ok {
			return strconv.FormatFloat(f, 'g', -1, 64)
		}
		return "0"
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
			parts[i] = goLiteral(deref(t.elem), e)
		}
		return goTypeStr(t) + "{" + strings.Join(parts, ", ") + "}"
	default: // any / interface{}
		return goAny(v)
	}
}

// goScalar renders an integer-ish value (int from YAML, float64 from JSON) as a
// Go int literal.
func goScalar(v any) string {
	switch x := v.(type) {
	case int:
		return strconv.Itoa(x)
	case int64:
		return strconv.FormatInt(x, 10)
	case int32:
		return strconv.FormatInt(int64(x), 10)
	case float64:
		return strconv.FormatInt(int64(x), 10)
	}
	return "0"
}

// asFloat coerces a JSON/YAML number to float64.
func asFloat(v any) (float64, bool) {
	switch x := v.(type) {
	case float64:
		return x, true
	case int:
		return float64(x), true
	case int64:
		return float64(x), true
	case int32:
		return float64(x), true
	}
	return 0, false
}

// goAny renders an interface{} literal (best-effort, for mixed/null data).
func goAny(v any) string {
	switch x := v.(type) {
	case nil:
		return "nil"
	case bool:
		return strconv.FormatBool(x)
	case int, int64, int32:
		return goScalar(x)
	case float64:
		if x == math.Trunc(x) {
			return strconv.FormatInt(int64(x), 10)
		}
		return strconv.FormatFloat(x, 'g', -1, 64)
	case string:
		return strconv.Quote(x)
	case []any:
		parts := make([]string, len(x))
		for i, e := range x {
			parts[i] = goAny(e)
		}
		return "[]interface{}{" + strings.Join(parts, ", ") + "}"
	default:
		return "nil"
	}
}

// nodeBuildLit emits buildList/buildTree(<array literal>) for a node argument.
func nodeBuildLit(shapeKind string, v any) string {
	arr, _ := v.([]any)
	if shapeKind == "tree" {
		parts := make([]string, len(arr))
		for i, e := range arr {
			if e == nil {
				parts[i] = "nil"
			} else {
				parts[i] = goScalar(e)
			}
		}
		return "buildTree([]interface{}{" + strings.Join(parts, ", ") + "})"
	}
	parts := make([]string, len(arr))
	for i, e := range arr {
		parts[i] = goScalar(e)
	}
	return "buildList([]int{" + strings.Join(parts, ", ") + "})"
}

// goPrelude returns the ListNode/TreeNode types + build/dump helpers (ported from
// the Python shapePrelude) for node problems.
func goPrelude(shapeKind string) string {
	if shapeKind == "linkedlist" {
		return `type ListNode struct {
	Val  int
	Next *ListNode
}

func buildList(a []int) *ListNode {
	var head *ListNode
	for i := len(a) - 1; i >= 0; i-- {
		head = &ListNode{Val: a[i], Next: head}
	}
	return head
}

func dumpList(n *ListNode) []int {
	out := []int{}
	for n != nil {
		out = append(out, n.Val)
		n = n.Next
	}
	return out
}

`
	}
	if shapeKind == "tree" {
		return `type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func buildTree(a []interface{}) *TreeNode {
	if len(a) == 0 || a[0] == nil {
		return nil
	}
	toInt := func(v interface{}) int {
		switch x := v.(type) {
		case int:
			return x
		case float64:
			return int(x)
		}
		return 0
	}
	root := &TreeNode{Val: toInt(a[0])}
	q := []*TreeNode{root}
	i := 1
	for len(q) > 0 && i < len(a) {
		node := q[0]
		q = q[1:]
		if i < len(a) {
			if a[i] != nil {
				node.Left = &TreeNode{Val: toInt(a[i])}
				q = append(q, node.Left)
			}
			i++
		}
		if i < len(a) {
			if a[i] != nil {
				node.Right = &TreeNode{Val: toInt(a[i])}
				q = append(q, node.Right)
			}
			i++
		}
	}
	return root
}

func dumpTree(root *TreeNode) []interface{} {
	out := []interface{}{}
	if root == nil {
		return out
	}
	q := []*TreeNode{root}
	for len(q) > 0 {
		node := q[0]
		q = q[1:]
		if node == nil {
			out = append(out, nil)
			continue
		}
		out = append(out, node.Val)
		q = append(q, node.Left, node.Right)
	}
	for len(out) > 0 && out[len(out)-1] == nil {
		out = out[:len(out)-1]
	}
	return out
}

`
	}
	return ""
}

// argTypesFor infers the Go type of each argument position from the test data
// (node args come from the Shape, not inference).
func argTypesFor(tests []TestCase, shape Shape) []gtype {
	n := 0
	if len(tests) > 0 {
		n = len(tests[0].Input)
	}
	out := make([]gtype, n)
	for i := 0; i < n; i++ {
		if i < len(shape.ArgKinds) && shape.ArgKinds[i] == "node" {
			out[i] = nodeType(shape.Kind)
			continue
		}
		var samples []any
		for _, tc := range tests {
			if i < len(tc.Input) {
				samples = append(samples, tc.Input[i])
			}
		}
		out[i] = inferType(samples)
	}
	return out
}

// BuildGoDriver assembles the full Go program: prelude (if node) + player source
// + a main that calls funcName on each test's embedded-literal args and prints
// the results in the shared marker format.
func BuildGoDriver(source, funcName string, tests []TestCase, shape Shape) (string, error) {
	argTypes := argTypesFor(tests, shape)
	var b strings.Builder
	b.WriteString("package main\n\nimport (\n\t\"encoding/json\"\n\t\"fmt\"\n)\n\n")
	if shape.Kind != "" {
		b.WriteString(goPrelude(shape.Kind))
	}
	b.WriteString(strings.TrimRight(source, "\n") + "\n\n")
	b.WriteString("func main() {\n")
	b.WriteString("\ttype __r struct {\n\t\tName  string      `json:\"name\"`\n\t\tGot   interface{} `json:\"got\"`\n\t\tError string      `json:\"error\"`\n\t}\n")
	b.WriteString("\t__out := []__r{}\n")
	for _, tc := range tests {
		argStrs := make([]string, len(argTypes))
		for i := range argTypes {
			var v any
			if i < len(tc.Input) {
				v = tc.Input[i]
			}
			if i < len(shape.ArgKinds) && shape.ArgKinds[i] == "node" {
				argStrs[i] = nodeBuildLit(shape.Kind, v)
			} else {
				argStrs[i] = goLiteral(argTypes[i], v)
			}
		}
		call := funcName + "(" + strings.Join(argStrs, ", ") + ")"
		if shape.RetKind == "node" {
			if shape.Kind == "tree" {
				call = "dumpTree(" + call + ")"
			} else {
				call = "dumpList(" + call + ")"
			}
		}
		q := strconv.Quote(tc.Name)
		b.WriteString("\tfunc() {\n")
		b.WriteString("\t\tdefer func() {\n\t\t\tif __e := recover(); __e != nil {\n\t\t\t\t__out = append(__out, __r{Name: " + q + ", Error: fmt.Sprint(__e)})\n\t\t\t}\n\t\t}()\n")
		b.WriteString("\t\t__out = append(__out, __r{Name: " + q + ", Got: " + call + "})\n")
		b.WriteString("\t}()\n")
	}
	b.WriteString("\t__b, _ := json.Marshal(__out)\n")
	b.WriteString("\tfmt.Println(" + strconv.Quote(marker) + " + string(__b))\n")
	b.WriteString("}\n")
	return b.String(), nil
}

var pyDefRe = regexp.MustCompile(`def\s+\w+\s*\(([^)]*)\)`)

// paramNames extracts parameter names from a Python def (for nicer Go stubs),
// falling back to a0, a1, … . self/cls and defaults/annotations are stripped.
func paramNames(pySource string, n int) []string {
	names := make([]string, n)
	for i := range names {
		names[i] = fmt.Sprintf("a%d", i)
	}
	if m := pyDefRe.FindStringSubmatch(pySource); m != nil {
		raw := strings.Split(m[1], ",")
		var parsed []string
		for _, p := range raw {
			p = strings.TrimSpace(p)
			if p == "" || p == "self" || p == "cls" {
				continue
			}
			if j := strings.IndexAny(p, ":="); j >= 0 {
				p = strings.TrimSpace(p[:j])
			}
			if p != "" {
				parsed = append(parsed, p)
			}
		}
		for i := 0; i < n && i < len(parsed); i++ {
			names[i] = parsed[i]
		}
	}
	return names
}

func goZero(t gtype) string {
	switch t.kind {
	case "int", "float":
		return "0"
	case "string":
		return `""`
	case "bool":
		return "false"
	case "slice":
		return "nil"
	case "list", "tree":
		return "nil"
	default:
		return "nil"
	}
}

// GoStarter generates a Go function stub for a problem, inferred from the test
// data (params named from the Python source where available). retKind=node →
// the player returns a *ListNode/*TreeNode (the harness dumps it).
func GoStarter(funcName, pySource string, tests []TestCase, shape Shape) string {
	argTypes := argTypesFor(tests, shape)
	names := paramNames(pySource, len(argTypes))
	params := make([]string, len(argTypes))
	for i, t := range argTypes {
		params[i] = names[i] + " " + goTypeStr(t)
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
	retStr := goTypeStr(ret)
	return fmt.Sprintf("func %s(%s) %s {\n\t// your code here\n\treturn %s\n}\n",
		funcName, strings.Join(params, ", "), retStr, goZero(ret))
}
