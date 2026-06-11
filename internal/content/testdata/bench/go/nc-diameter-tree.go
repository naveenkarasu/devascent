func diameter_of_binary_tree(root *TreeNode) int {
	best := 0
	var height func(node *TreeNode) int
	height = func(node *TreeNode) int {
		if node == nil {
			return 0
		}
		l := height(node.Left)
		r := height(node.Right)
		if l+r > best {
			best = l + r
		}
		if l > r {
			return 1 + l
		}
		return 1 + r
	}
	height(root)
	return best
}
