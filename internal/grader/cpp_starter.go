package grader

// C++ is reference/view-only (no grader harness yet), but every surface that
// seeds an editor still needs an idiomatic typed stub — serving the authored
// Python starter to a C++ session is the same release-blocking bug class the
// GUI entrance test had. Only the starter generator lives here; grading
// availability is still gated in engine.GradingAvailable.

import (
	"fmt"
	"strings"
)

func cppTypeStr(t gtype) string {
	switch t.kind {
	case "int":
		return "long long"
	case "float":
		return "double"
	case "string":
		return "std::string"
	case "bool":
		return "bool"
	case "slice":
		return "std::vector<" + cppTypeStr(deref(t.elem)) + ">"
	case "list":
		return "ListNode*"
	case "tree":
		return "TreeNode*"
	default:
		return "std::any"
	}
}

// CppStarter generates a view-only C++ stub ({} value-initializes every
// return type, so the stub is valid C++ as written).
func CppStarter(funcName, pySource string, tests []TestCase, shape Shape) string {
	argTypes := argTypesFor(tests, shape)
	names := paramNames(pySource, len(argTypes))
	params := make([]string, len(argTypes))
	for i, t := range argTypes {
		params[i] = cppTypeStr(t) + " " + names[i]
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
	sig := fmt.Sprintf("%s %s(%s)", cppTypeStr(ret), funcName, strings.Join(params, ", "))

	var inc []string
	for _, h := range []string{"string", "vector", "any"} {
		if strings.Contains(sig, "std::"+h) {
			inc = append(inc, "#include <"+h+">")
		}
	}
	head := ""
	if len(inc) > 0 {
		head = strings.Join(inc, "\n") + "\n\n"
	}
	return fmt.Sprintf("%s%s {\n    // your code here\n    return {};\n}\n", head, sig)
}
