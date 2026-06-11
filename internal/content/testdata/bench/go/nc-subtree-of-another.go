func is_subtree(root *TreeNode, sub *TreeNode) bool {
	var same func(a, b *TreeNode) bool
	same = func(a, b *TreeNode) bool {
		if a == nil && b == nil {
			return true
		}
		if a == nil || b == nil || a.Val != b.Val {
			return false
		}
		return same(a.Left, b.Left) && same(a.Right, b.Right)
	}
	if sub == nil {
		return true
	}
	var dfs func(node *TreeNode) bool
	dfs = func(node *TreeNode) bool {
		if node == nil {
			return false
		}
		if same(node, sub) {
			return true
		}
		return dfs(node.Left) || dfs(node.Right)
	}
	return dfs(root)
}
