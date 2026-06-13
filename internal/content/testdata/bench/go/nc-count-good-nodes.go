func good_nodes(root *TreeNode) int {
	var dfs func(node *TreeNode, mx int) int
	dfs = func(node *TreeNode, mx int) int {
		if node == nil {
			return 0
		}
		count := 0
		if node.Val >= mx {
			count = 1
		}
		if node.Val > mx {
			mx = node.Val
		}
		return count + dfs(node.Left, mx) + dfs(node.Right, mx)
	}
	return dfs(root, -1<<62)
}
