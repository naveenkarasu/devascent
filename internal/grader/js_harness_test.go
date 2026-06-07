package grader

import (
	"testing"
	"time"

	"devascent/internal/toolchain"
)

func jsTSCases() []struct {
	name, fn, src string
	tests         []TestCase
	shape         Shape
} {
	return []struct {
		name, fn, src string
		tests         []TestCase
		shape         Shape
	}{
		{"scalar add", "add", "function add(a, b) { return a + b; }",
			[]TestCase{{Name: "t1", Input: []any{2, 3}, Expected: 5}, {Name: "t2", Input: []any{-1, 1}, Expected: 0}}, Shape{}},
		{"slice arg + ret (two-sum)", "twoSum",
			"function twoSum(nums, target) { const seen = {}; for (let i = 0; i < nums.length; i++) { if ((target - nums[i]) in seen) return [seen[target - nums[i]], i]; seen[nums[i]] = i; } return []; }",
			[]TestCase{{Name: "t1", Input: []any{[]any{2, 7, 11, 15}, 9}, Expected: []any{0, 1}}}, Shape{}},
		{"string", "reverseStr", "function reverseStr(s) { return s.split('').reverse().join(''); }",
			[]TestCase{{Name: "t1", Input: []any{"abc"}, Expected: "cba"}}, Shape{}},
		{"bool", "isEven", "function isEven(n) { return n % 2 === 0; }",
			[]TestCase{{Name: "t1", Input: []any{4}, Expected: true}, {Name: "t2", Input: []any{3}, Expected: false}}, Shape{}},
		{"linked-list (node arg + ret)", "reverseList",
			"function reverseList(head) { let prev = null; while (head) { const next = head.next; head.next = prev; prev = head; head = next; } return prev; }",
			[]TestCase{{Name: "t1", Input: []any{[]any{1, 2, 3}}, Expected: []any{3, 2, 1}}},
			Shape{Kind: "linkedlist", ArgKinds: []string{"node"}, RetKind: "node"}},
		{"tree (node arg, scalar ret)", "maxDepth",
			"function maxDepth(root) { if (!root) return 0; return Math.max(maxDepth(root.left), maxDepth(root.right)) + 1; }",
			[]TestCase{{Name: "t1", Input: []any{[]any{3, 9, 20, nil, nil, 15, 7}}, Expected: 3}},
			Shape{Kind: "tree", ArgKinds: []string{"node"}}},
	}
}

func TestJSHarness_TypeShapes(t *testing.T) {
	if testing.Short() {
		t.Skip("node grading is slow-ish; skipped in -short")
	}
	g := NewLocalToolchain(toolchain.New())
	if g.det.Presence("javascript").Status == toolchain.Missing {
		t.Skip("node not installed")
	}
	g.timeout = 30 * time.Second
	for _, c := range jsTSCases() {
		t.Run(c.name, func(t *testing.T) {
			v, err := g.Run("javascript", c.src, c.fn, c.tests, c.shape)
			if err != nil {
				t.Fatal(err)
			}
			if !v.Passed {
				t.Fatalf("%s: should pass, got err=%q results=%+v", c.name, v.Err, v.Results)
			}
		})
	}
	v, _ := g.Run("javascript", "function add(a, b) { return a - b; }", "add",
		[]TestCase{{Name: "t1", Input: []any{2, 3}, Expected: 5}}, Shape{})
	if v.Passed {
		t.Fatal("a-b should not pass")
	}
}

func TestTSHarness_TypeShapes(t *testing.T) {
	if testing.Short() {
		t.Skip("TS compile is slow; skipped in -short")
	}
	g := NewLocalToolchain(toolchain.New())
	if g.det.Presence("typescript").Status == toolchain.Missing {
		t.Skip("tsc/node not installed")
	}
	g.timeout = 60 * time.Second
	// Reuse the JS sources — they are valid TS (untyped → implicit any compiles
	// under default tsc; the call is cast `as any` by the harness).
	for _, c := range jsTSCases() {
		t.Run(c.name, func(t *testing.T) {
			v, err := g.Run("typescript", c.src, c.fn, c.tests, c.shape)
			if err != nil {
				t.Fatal(err)
			}
			if !v.Passed {
				t.Fatalf("%s: should pass, got err=%q results=%+v", c.name, v.Err, v.Results)
			}
		})
	}
}
