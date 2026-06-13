func is_valid_bst(root *TreeNode) bool {
	var valid func(node *TreeNode, lo, hi int) bool
	valid = func(node *TreeNode, lo, hi int) bool {
		if node == nil {
			return true
		}
		if node.Val <= lo || node.Val >= hi {
			return false
		}
		return valid(node.Left, lo, node.Val) && valid(node.Right, node.Val, hi)
	}
	return valid(root, -1<<62, 1<<62)
}
