package grader

import (
	"testing"
	"time"

	"devascent/internal/toolchain"
)

func TestJavaHarness_TypeShapes(t *testing.T) {
	if testing.Short() {
		t.Skip("java grading is slow; skipped in -short")
	}
	g := NewLocalToolchain(toolchain.New())
	if g.det.Presence("java").Status == toolchain.Missing {
		t.Skip("JDK not installed")
	}
	g.timeout = 90 * time.Second
	type tc struct {
		name, fn, src string
		tests         []TestCase
		shape         Shape
	}
	cases := []tc{
		{"scalar add", "add", "class Solution { public long add(long a, long b){ return a + b; } }",
			[]TestCase{{Name: "t1", Input: []any{2, 3}, Expected: 5}, {Name: "t2", Input: []any{-1, 1}, Expected: 0}}, Shape{}},
		{"slice arg + ret (two-sum)", "twoSum",
			"import java.util.*; class Solution { public long[] twoSum(long[] nums, long target){ Map<Long,Integer> m = new HashMap<>(); for (int i = 0; i < nums.length; i++){ if (m.containsKey(target - nums[i])) return new long[]{ m.get(target - nums[i]), i }; m.put(nums[i], i); } return new long[0]; } }",
			[]TestCase{{Name: "t1", Input: []any{[]any{2, 7, 11, 15}, 9}, Expected: []any{0, 1}}}, Shape{}},
		{"string", "reverseStr", "class Solution { public String reverseStr(String s){ return new StringBuilder(s).reverse().toString(); } }",
			[]TestCase{{Name: "t1", Input: []any{"abc"}, Expected: "cba"}}, Shape{}},
		{"bool", "isEven", "class Solution { public boolean isEven(long n){ return n % 2 == 0; } }",
			[]TestCase{{Name: "t1", Input: []any{4}, Expected: true}, {Name: "t2", Input: []any{3}, Expected: false}}, Shape{}},
		{"linked-list (node arg + ret)", "reverseList",
			"class Solution { public ListNode reverseList(ListNode head){ ListNode prev = null; while (head != null){ ListNode next = head.next; head.next = prev; prev = head; head = next; } return prev; } }",
			[]TestCase{{Name: "t1", Input: []any{[]any{1, 2, 3}}, Expected: []any{3, 2, 1}}}, Shape{Kind: "linkedlist", ArgKinds: []string{"node"}, RetKind: "node"}},
		{"tree (node arg, scalar ret)", "maxDepth",
			"class Solution { public long maxDepth(TreeNode root){ if (root == null) return 0; return Math.max(maxDepth(root.left), maxDepth(root.right)) + 1; } }",
			[]TestCase{{Name: "t1", Input: []any{[]any{3, 9, 20, nil, nil, 15, 7}}, Expected: 3}}, Shape{Kind: "tree", ArgKinds: []string{"node"}}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			v, err := g.Run("java", c.src, c.fn, c.tests, c.shape)
			if err != nil {
				t.Fatal(err)
			}
			if !v.Passed {
				t.Fatalf("%s: should pass, got err=%q results=%+v", c.name, v.Err, v.Results)
			}
		})
	}
	v, _ := g.Run("java", "class Solution { public long add(long a, long b){ return a - b; } }", "add",
		[]TestCase{{Name: "t1", Input: []any{2, 3}, Expected: 5}}, Shape{})
	if v.Passed {
		t.Fatal("a-b should not pass")
	}
}

func TestRustHarness_TypeShapes(t *testing.T) {
	if testing.Short() {
		t.Skip("rust grading is slow; skipped in -short")
	}
	g := NewLocalToolchain(toolchain.New())
	if g.det.Presence("rust").Status == toolchain.Missing {
		t.Skip("rust not installed")
	}
	g.timeout = 90 * time.Second
	type tc struct {
		name, fn, src string
		tests         []TestCase
		shape         Shape
	}
	cases := []tc{
		{"scalar add", "add", "fn add(a: i64, b: i64) -> i64 { a + b }",
			[]TestCase{{Name: "t1", Input: []any{2, 3}, Expected: 5}, {Name: "t2", Input: []any{-1, 1}, Expected: 0}}, Shape{}},
		{"slice arg + ret (two-sum)", "two_sum",
			"use std::collections::HashMap;\nfn two_sum(nums: Vec<i64>, target: i64) -> Vec<i64> { let mut seen: HashMap<i64, i64> = HashMap::new(); for (i, &n) in nums.iter().enumerate() { if let Some(&j) = seen.get(&(target - n)) { return vec![j, i as i64]; } seen.insert(n, i as i64); } vec![] }",
			[]TestCase{{Name: "t1", Input: []any{[]any{2, 7, 11, 15}, 9}, Expected: []any{0, 1}}}, Shape{}},
		{"string", "reverse_str", "fn reverse_str(s: String) -> String { s.chars().rev().collect() }",
			[]TestCase{{Name: "t1", Input: []any{"abc"}, Expected: "cba"}}, Shape{}},
		{"bool", "is_even", "fn is_even(n: i64) -> bool { n % 2 == 0 }",
			[]TestCase{{Name: "t1", Input: []any{4}, Expected: true}, {Name: "t2", Input: []any{3}, Expected: false}}, Shape{}},
		{"linked-list (node arg + ret)", "reverse_list",
			"fn reverse_list(mut head: Option<Box<ListNode>>) -> Option<Box<ListNode>> { let mut prev = None; while let Some(mut node) = head { head = node.next.take(); node.next = prev; prev = Some(node); } prev }",
			[]TestCase{{Name: "t1", Input: []any{[]any{1, 2, 3}}, Expected: []any{3, 2, 1}}}, Shape{Kind: "linkedlist", ArgKinds: []string{"node"}, RetKind: "node"}},
		{"tree (node arg, scalar ret)", "max_depth",
			"use std::rc::Rc;\nuse std::cell::RefCell;\nfn max_depth(root: Option<Rc<RefCell<TreeNode>>>) -> i64 { match root { None => 0, Some(n) => { let l = max_depth(n.borrow().left.clone()); let r = max_depth(n.borrow().right.clone()); 1 + l.max(r) } } }",
			[]TestCase{{Name: "t1", Input: []any{[]any{3, 9, 20, nil, nil, 15, 7}}, Expected: 3}}, Shape{Kind: "tree", ArgKinds: []string{"node"}}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			v, err := g.Run("rust", c.src, c.fn, c.tests, c.shape)
			if err != nil {
				t.Fatal(err)
			}
			if !v.Passed {
				t.Fatalf("%s: should pass, got err=%q results=%+v", c.name, v.Err, v.Results)
			}
		})
	}
	v, _ := g.Run("rust", "fn add(a: i64, b: i64) -> i64 { a - b }", "add",
		[]TestCase{{Name: "t1", Input: []any{2, 3}, Expected: 5}}, Shape{})
	if v.Passed {
		t.Fatal("a-b should not pass")
	}
}
