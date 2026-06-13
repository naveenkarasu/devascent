package grader

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// rust_harness.go — Model A function-call grading for Rust via IN-LANGUAGE
// comparison (no stdlib JSON). The player writes `fn func_name(...) -> Ret {...}`;
// the harness embeds args + expected as Rust literals, compares with == (after
// dumping nodes to Vec), catches panics per case (catch_unwind), and prints the
// in-language line protocol (ParseInLangOutput).

func rustTypeStr(t gtype) string {
	switch t.kind {
	case "int":
		return "i64"
	case "float":
		return "f64"
	case "string":
		return "String"
	case "bool":
		return "bool"
	case "slice":
		return "Vec<" + rustTypeStr(deref(t.elem)) + ">"
	case "list":
		return "Option<Box<ListNode>>"
	case "tree":
		return "Option<Rc<RefCell<TreeNode>>>"
	case "any":
		return "J"
	default:
		return "()"
	}
}

// rustJsonLit renders a heterogeneous (interface{}) value as a J literal — the
// analogue of goAny. Rust has no native `any`, so the harness provides a JSON
// value enum `J`; op-list/design inputs and mixed/null/dict returns use it.
// Integral floats become J::Int (matching the player's integral-vs-fractional
// branch); maps become J::Obj with SORTED keys (== on Vec is order-sensitive).
func rustJsonLit(v any) string {
	switch x := v.(type) {
	case nil:
		return "J::Null"
	case bool:
		return "J::Bool(" + strconv.FormatBool(x) + ")"
	case int, int64, int32:
		return "J::Int(" + goScalar(x) + ")"
	case float64:
		if x == float64(int64(x)) {
			return "J::Int(" + strconv.FormatInt(int64(x), 10) + ")"
		}
		s := strconv.FormatFloat(x, 'g', -1, 64)
		if !strings.ContainsAny(s, ".eE") {
			s += ".0"
		}
		return "J::Flt(" + s + ")"
	case string:
		return "J::Str(" + strconv.Quote(x) + ".to_string())"
	case []any:
		parts := make([]string, len(x))
		for i, e := range x {
			parts[i] = rustJsonLit(e)
		}
		return "J::Arr(vec![" + strings.Join(parts, ", ") + "])"
	case map[string]any:
		keys := make([]string, 0, len(x))
		for k := range x {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		parts := make([]string, len(keys))
		for i, k := range keys {
			parts[i] = "(" + strconv.Quote(k) + ".to_string(), " + rustJsonLit(x[k]) + ")"
		}
		return "J::Obj(vec![" + strings.Join(parts, ", ") + "])"
	default:
		return "J::Null"
	}
}

func rustLiteral(t gtype, v any) string {
	switch t.kind {
	case "int":
		return goScalar(v) + "i64" // explicit width so Vec<i64> comparisons type-check
	case "float":
		if f, ok := asFloat(v); ok {
			s := strconv.FormatFloat(f, 'g', -1, 64)
			if !strings.ContainsAny(s, ".eE") {
				s += ".0"
			}
			return s
		}
		return "0.0"
	case "string":
		s, _ := v.(string)
		return strconv.Quote(s) + ".to_string()"
	case "bool":
		b, _ := v.(bool)
		return strconv.FormatBool(b)
	case "slice":
		arr, _ := v.([]any)
		if len(arr) == 0 {
			// A bare vec![] gives Rust no element type to infer (E0282) when the
			// literal stands alone, e.g. `let __exp = vec![]` for a Vec<String>
			// return. Emit a typed empty vec so any element type round-trips.
			return "Vec::<" + rustTypeStr(deref(t.elem)) + ">::new()"
		}
		parts := make([]string, len(arr))
		for i, e := range arr {
			parts[i] = rustLiteral(deref(t.elem), e)
		}
		return "vec![" + strings.Join(parts, ", ") + "]"
	case "any":
		return rustJsonLit(v)
	default:
		return "()"
	}
}

func rustNodeBuildLit(shapeKind string, v any) string {
	arr, _ := v.([]any)
	if shapeKind == "tree" {
		parts := make([]string, len(arr))
		for i, e := range arr {
			if e == nil {
				parts[i] = "None"
			} else {
				parts[i] = "Some(" + goScalar(e) + ")"
			}
		}
		return "__build_tree(vec![" + strings.Join(parts, ", ") + "])"
	}
	parts := make([]string, len(arr))
	for i, e := range arr {
		parts[i] = goScalar(e)
	}
	return "__build_list(vec![" + strings.Join(parts, ", ") + "])"
}

func rustPrelude(shape Shape) string {
	var b strings.Builder
	// J: a JSON value enum so heterogeneous (`any`) op-list/mixed/dict data has a
	// Rust type with == (Rust has no native `any`; serde isn't available to rustc).
	// PartialEq only (f64 isn't Eq). Harness-provided — players use it, don't define it.
	b.WriteString("#[derive(Clone, Debug, PartialEq)]\nenum J { Null, Int(i64), Flt(f64), Str(String), Bool(bool), Arr(Vec<J>), Obj(Vec<(String, J)>) }\n")
	if shape.Kind == "linkedlist" {
		b.WriteString("#[derive(PartialEq, Eq, Clone, Debug)]\npub struct ListNode { pub val: i32, pub next: Option<Box<ListNode>> }\n")
		b.WriteString("impl ListNode { fn new(v: i32) -> Self { ListNode { val: v, next: None } } }\n")
		b.WriteString("fn __build_list(a: Vec<i32>) -> Option<Box<ListNode>> { let mut head = None; for &v in a.iter().rev() { let mut n = Box::new(ListNode::new(v)); n.next = head; head = Some(n); } head }\n")
		if shape.RetKind == "node" {
			b.WriteString("fn __dump_list(mut n: Option<Box<ListNode>>) -> Vec<i64> { let mut o = Vec::new(); while let Some(node) = n { o.push(node.val as i64); n = node.next; } o }\n")
		}
	}
	if shape.Kind == "tree" {
		// Fully-qualified std::rc::Rc / std::cell::RefCell (NO `use`) so the player's
		// own `use std::rc::Rc;` etc. don't collide (E0252). The starter supplies the
		// `use` lines for the player's editable code.
		b.WriteString("#[derive(PartialEq, Eq, Clone, Debug)]\npub struct TreeNode { pub val: i32, pub left: Option<std::rc::Rc<std::cell::RefCell<TreeNode>>>, pub right: Option<std::rc::Rc<std::cell::RefCell<TreeNode>>> }\n")
		b.WriteString("impl TreeNode { fn new(v: i32) -> Self { TreeNode { val: v, left: None, right: None } } }\n")
		b.WriteString("fn __build_tree(a: Vec<Option<i32>>) -> Option<std::rc::Rc<std::cell::RefCell<TreeNode>>> { if a.is_empty() || a[0].is_none() { return None; } let root = std::rc::Rc::new(std::cell::RefCell::new(TreeNode::new(a[0].unwrap()))); let mut q = std::collections::VecDeque::new(); q.push_back(std::rc::Rc::clone(&root)); let mut i = 1usize; while !q.is_empty() && i < a.len() { let node = q.pop_front().unwrap(); if i < a.len() { if let Some(v) = a[i] { let c = std::rc::Rc::new(std::cell::RefCell::new(TreeNode::new(v))); node.borrow_mut().left = Some(std::rc::Rc::clone(&c)); q.push_back(c); } i += 1; } if i < a.len() { if let Some(v) = a[i] { let c = std::rc::Rc::new(std::cell::RefCell::new(TreeNode::new(v))); node.borrow_mut().right = Some(std::rc::Rc::clone(&c)); q.push_back(c); } i += 1; } } Some(root) }\n")
		if shape.RetKind == "node" {
			b.WriteString("fn __dump_tree(root: Option<std::rc::Rc<std::cell::RefCell<TreeNode>>>) -> Vec<Option<i64>> { let mut o: Vec<Option<i64>> = Vec::new(); if root.is_none() { return o; } let mut q: std::collections::VecDeque<Option<std::rc::Rc<std::cell::RefCell<TreeNode>>>> = std::collections::VecDeque::new(); q.push_back(root); while let Some(cur) = q.pop_front() { match cur { None => o.push(None), Some(n) => { o.push(Some(n.borrow().val as i64)); q.push_back(n.borrow().left.clone()); q.push_back(n.borrow().right.clone()); } } } while let Some(last) = o.last() { if last.is_none() { o.pop(); } else { break; } } o }\n")
		}
	}
	if b.Len() > 0 {
		b.WriteString("\n")
	}
	return b.String()
}

func rustRetInfo(tests []TestCase, shape Shape) (ret gtype, isList, isTree bool) {
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

// rustExpectedLit renders the expected value as a Rust literal of the comparison
// type (node returns dump to Vec<i64> / Vec<Option<i64>>).
func rustExpectedLit(ret gtype, v any, isList, isTree bool) string {
	if isList {
		arr, _ := v.([]any)
		if len(arr) == 0 {
			return "Vec::<i64>::new()"
		}
		parts := make([]string, len(arr))
		for i, e := range arr {
			parts[i] = goScalar(e) + "i64"
		}
		return "vec![" + strings.Join(parts, ", ") + "]"
	}
	if isTree {
		arr, _ := v.([]any)
		if len(arr) == 0 {
			return "Vec::<Option<i64>>::new()"
		}
		parts := make([]string, len(arr))
		for i, e := range arr {
			if e == nil {
				parts[i] = "None"
			} else {
				parts[i] = "Some(" + goScalar(e) + "i64)"
			}
		}
		return "vec![" + strings.Join(parts, ", ") + "]"
	}
	return rustLiteral(ret, v)
}

// rustCompareExpr returns the boolean comparison of __got vs __exp (float uses a
// tolerance to match jsonEqual; everything else is exact ==).
func rustCompareExpr(ret gtype, isList, isTree bool) string {
	if !isList && !isTree && ret.kind == "float" {
		return "(__got - __exp).abs() < 1e-6"
	}
	return "__got == __exp"
}

// BuildRustDriver assembles the full Rust program.
func BuildRustDriver(source, funcName string, tests []TestCase, shape Shape) (string, error) {
	argTypes := argTypesFor(tests, shape)
	ret, isList, isTree := rustRetInfo(tests, shape)
	var b strings.Builder
	b.WriteString("#![allow(warnings)]\n")
	b.WriteString(rustPrelude(shape))
	b.WriteString(strings.TrimRight(source, "\n") + "\n\n")
	b.WriteString("fn main() {\n")
	b.WriteString("    std::panic::set_hook(Box::new(|_| {}));\n")
	b.WriteString("    const M: &str = " + strconv.Quote(marker) + ";\n")
	for _, tc := range tests {
		args := make([]string, len(argTypes))
		for i := range argTypes {
			var v any
			if i < len(tc.Input) {
				v = tc.Input[i]
			}
			if i < len(shape.ArgKinds) && shape.ArgKinds[i] == "node" {
				args[i] = rustNodeBuildLit(shape.Kind, v)
			} else {
				args[i] = rustLiteral(argTypes[i], v)
			}
		}
		call := funcName + "(" + strings.Join(args, ", ") + ")"
		if isList {
			call = "__dump_list(" + call + ")"
		} else if isTree {
			call = "__dump_tree(" + call + ")"
		}
		exp := rustExpectedLit(ret, tc.Expected, isList, isTree)
		q := strconv.Quote(tc.Name)
		b.WriteString("    {\n")
		b.WriteString("        let __r = std::panic::catch_unwind(|| " + call + ");\n")
		b.WriteString("        match __r {\n")
		b.WriteString("            Ok(__got) => { let __exp = " + exp + "; let __ok = " + rustCompareExpr(ret, isList, isTree) + "; println!(\"{}{}\\u{1f}{}\\u{1f}{:?}\", M, " + q + ", if __ok {\"1\"} else {\"0\"}, __got); }\n")
		b.WriteString("            Err(_) => println!(\"{}{}\\u{1f}0\\u{1f}panic\", M, " + q + "),\n")
		b.WriteString("        }\n")
		b.WriteString("    }\n")
	}
	b.WriteString("}\n")
	return b.String(), nil
}

func rustZero(t gtype) string {
	switch t.kind {
	case "int":
		return "0"
	case "float":
		return "0.0"
	case "string":
		return "String::new()"
	case "bool":
		return "false"
	case "slice":
		return "vec![]"
	case "list", "tree":
		return "None"
	default:
		return "()"
	}
}

// RustStarter generates a Rust function stub for a problem.
func RustStarter(funcName, pySource string, tests []TestCase, shape Shape) string {
	argTypes := argTypesFor(tests, shape)
	names := paramNames(pySource, len(argTypes))
	params := make([]string, len(argTypes))
	for i, t := range argTypes {
		params[i] = names[i] + ": " + rustTypeStr(t)
	}
	var ret gtype
	if shape.RetKind == "node" {
		ret = nodeType(shape.Kind)
	} else {
		ret, _, _ = rustRetInfo(tests, shape)
	}
	stub := fmt.Sprintf("fn %s(%s) -> %s {\n    // your code here\n    %s\n}\n",
		funcName, strings.Join(params, ", "), rustTypeStr(ret), rustZero(ret))
	if shape.Kind == "tree" { // tree problems need Rc/RefCell in scope for the player
		stub = "use std::rc::Rc;\nuse std::cell::RefCell;\n\n" + stub
	}
	return stub
}
