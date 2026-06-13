package grader

import (
	"testing"

	"devascent/internal/toolchain"
)

func goGrader(t *testing.T) *LocalToolchain {
	t.Helper()
	g := NewLocalToolchain(toolchain.New())
	if g.det.Presence("go").Status == toolchain.Missing {
		t.Skip("go toolchain not installed")
	}
	g.timeout = 60 * 1000 * 1000 * 1000 // 60s (go run can be cold)
	return g
}

// TestGoHarness_TypeShapes grades real Go solutions across the type shapes the
// inference must handle: scalars, slices (arg + return), strings, bools, and the
// linked-list / tree node harnesses. This validates the function-call mechanism
// end-to-end (Model A) through the real Go toolchain.
func TestGoHarness_TypeShapes(t *testing.T) {
	g := goGrader(t)
	type tc struct {
		name   string
		fn     string
		src    string
		tests  []TestCase
		shape  Shape
		expect bool
	}
	cases := []tc{
		{
			name: "scalar add", fn: "add",
			src:    "func add(a int, b int) int { return a + b }",
			tests:  []TestCase{{Name: "t1", Input: []any{2, 3}, Expected: 5}, {Name: "t2", Input: []any{-1, 1}, Expected: 0}},
			expect: true,
		},
		{
			name: "slice arg, scalar ret", fn: "sumList",
			src:    "func sumList(nums []int) int {\n\ts := 0\n\tfor _, n := range nums {\n\t\ts += n\n\t}\n\treturn s\n}",
			tests:  []TestCase{{Name: "t1", Input: []any{[]any{1, 2, 3, 4}}, Expected: 10}},
			expect: true,
		},
		{
			name: "slice arg + ret (two-sum)", fn: "twoSum",
			src:    "func twoSum(nums []int, target int) []int {\n\tseen := map[int]int{}\n\tfor i, n := range nums {\n\t\tif j, ok := seen[target-n]; ok {\n\t\t\treturn []int{j, i}\n\t\t}\n\t\tseen[n] = i\n\t}\n\treturn nil\n}",
			tests:  []TestCase{{Name: "t1", Input: []any{[]any{2, 7, 11, 15}, 9}, Expected: []any{0, 1}}},
			expect: true,
		},
		{
			name: "string", fn: "reverseStr",
			src:    "func reverseStr(s string) string {\n\tr := []rune(s)\n\tfor i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {\n\t\tr[i], r[j] = r[j], r[i]\n\t}\n\treturn string(r)\n}",
			tests:  []TestCase{{Name: "t1", Input: []any{"abc"}, Expected: "cba"}},
			expect: true,
		},
		{
			name: "bool", fn: "isEven",
			src:    "func isEven(n int) bool { return n%2 == 0 }",
			tests:  []TestCase{{Name: "t1", Input: []any{4}, Expected: true}, {Name: "t2", Input: []any{3}, Expected: false}},
			expect: true,
		},
		{
			name: "linked-list (node arg + ret)", fn: "reverseList",
			src:   "func reverseList(head *ListNode) *ListNode {\n\tvar prev *ListNode\n\tfor head != nil {\n\t\tnext := head.Next\n\t\thead.Next = prev\n\t\tprev = head\n\t\thead = next\n\t}\n\treturn prev\n}",
			tests: []TestCase{{Name: "t1", Input: []any{[]any{1, 2, 3}}, Expected: []any{3, 2, 1}}},
			shape: Shape{Kind: "linkedlist", ArgKinds: []string{"node"}, RetKind: "node"}, expect: true,
		},
		{
			name: "tree (node arg, scalar ret)", fn: "maxDepth",
			src:   "func maxDepth(root *TreeNode) int {\n\tif root == nil {\n\t\treturn 0\n\t}\n\tl := maxDepth(root.Left)\n\tr := maxDepth(root.Right)\n\tif l > r {\n\t\treturn l + 1\n\t}\n\treturn r + 1\n}",
			tests: []TestCase{{Name: "t1", Input: []any{[]any{3, 9, 20, nil, nil, 15, 7}}, Expected: 3}},
			shape: Shape{Kind: "tree", ArgKinds: []string{"node"}}, expect: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			v, err := g.Run("go", c.src, c.fn, c.tests, c.shape)
			if err != nil {
				t.Fatalf("grader error: %v", err)
			}
			if v.Passed != c.expect {
				t.Fatalf("%s: passed=%v want %v (err=%q results=%+v)", c.name, v.Passed, c.expect, v.Err, v.Results)
			}
		})
	}

	// A wrong solution must FAIL (guards against a no-op grader).
	v, _ := g.Run("go", "func add(a int, b int) int { return a - b }", "add",
		[]TestCase{{Name: "t1", Input: []any{2, 3}, Expected: 5}}, Shape{})
	if v.Passed {
		t.Fatal("a-b should not pass the add tests")
	}
}
