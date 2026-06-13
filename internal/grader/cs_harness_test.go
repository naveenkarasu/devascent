package grader

import (
	"testing"
	"time"

	"devascent/internal/toolchain"
)

// TestCSharpHarness_TypeShapes grades real C# solutions across the type shapes
// the inference must handle, through the real .NET SDK. dotnet is slow (cold
// restore), so this is -short-skipped and uses a focused set.
func TestCSharpHarness_TypeShapes(t *testing.T) {
	if testing.Short() {
		t.Skip("C# grading via dotnet is slow; skipped in -short")
	}
	g := NewLocalToolchain(toolchain.New())
	if g.det.Presence("csharp").Status == toolchain.Missing {
		t.Skip(".NET SDK not installed")
	}
	g.timeout = 180 * time.Second // cold dotnet restore + build

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
			src:    "public class Solution { public long add(long a, long b){ return a + b; } }",
			tests:  []TestCase{{Name: "t1", Input: []any{2, 3}, Expected: 5}, {Name: "t2", Input: []any{-1, 1}, Expected: 0}},
			expect: true,
		},
		{
			name: "slice arg + ret (two-sum)", fn: "twoSum",
			src:    "using System.Collections.Generic;\npublic class Solution { public long[] twoSum(long[] nums, long target){ var seen = new Dictionary<long,int>(); for (int i = 0; i < nums.Length; i++){ if (seen.ContainsKey(target - nums[i])) return new long[]{ seen[target - nums[i]], i }; seen[nums[i]] = i; } return new long[0]; } }",
			tests:  []TestCase{{Name: "t1", Input: []any{[]any{2, 7, 11, 15}, 9}, Expected: []any{0, 1}}},
			expect: true,
		},
		{
			name: "string", fn: "reverseStr",
			src:    "public class Solution { public string reverseStr(string s){ var a = s.ToCharArray(); System.Array.Reverse(a); return new string(a); } }",
			tests:  []TestCase{{Name: "t1", Input: []any{"abc"}, Expected: "cba"}},
			expect: true,
		},
		{
			name: "bool", fn: "isEven",
			src:    "public class Solution { public bool isEven(long n){ return n % 2 == 0; } }",
			tests:  []TestCase{{Name: "t1", Input: []any{4}, Expected: true}, {Name: "t2", Input: []any{3}, Expected: false}},
			expect: true,
		},
		{
			name: "linked-list (node arg + ret)", fn: "reverseList",
			src:    "public class Solution { public ListNode reverseList(ListNode head){ ListNode prev = null; while (head != null){ var next = head.next; head.next = prev; prev = head; head = next; } return prev; } }",
			tests:  []TestCase{{Name: "t1", Input: []any{[]any{1, 2, 3}}, Expected: []any{3, 2, 1}}},
			shape:  Shape{Kind: "linkedlist", ArgKinds: []string{"node"}, RetKind: "node"},
			expect: true,
		},
		{
			name: "tree (node arg, scalar ret)", fn: "maxDepth",
			src:    "public class Solution { public long maxDepth(TreeNode root){ if (root == null) return 0; var l = maxDepth(root.left); var r = maxDepth(root.right); return (l > r ? l : r) + 1; } }",
			tests:  []TestCase{{Name: "t1", Input: []any{[]any{3, 9, 20, nil, nil, 15, 7}}, Expected: 3}},
			shape:  Shape{Kind: "tree", ArgKinds: []string{"node"}},
			expect: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			v, err := g.Run("csharp", c.src, c.fn, c.tests, c.shape)
			if err != nil {
				t.Fatalf("grader error: %v", err)
			}
			if v.Passed != c.expect {
				t.Fatalf("%s: passed=%v want %v (err=%q results=%+v)", c.name, v.Passed, c.expect, v.Err, v.Results)
			}
		})
	}

	v, _ := g.Run("csharp", "public class Solution { public long add(long a, long b){ return a - b; } }", "add",
		[]TestCase{{Name: "t1", Input: []any{2, 3}, Expected: 5}}, Shape{})
	if v.Passed {
		t.Fatal("a-b should not pass the add tests")
	}
}
