func max_depth(root *TreeNode) int {
	if root == nil {
		return 0
	}
	left := max_depth(root.Left)
	right := max_depth(root.Right)
	if left > right {
		return 1 + left
	}
	return 1 + right
}
