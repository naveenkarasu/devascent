func is_balanced(root *TreeNode) bool {
	var height func(node *TreeNode) int
	height = func(node *TreeNode) int {
		if node == nil {
			return 0
		}
		l := height(node.Left)
		if l < 0 {
			return -1
		}
		r := height(node.Right)
		if r < 0 {
			return -1
		}
		diff := l - r
		if diff < -1 || diff > 1 {
			return -1
		}
		if l > r {
			return 1 + l
		}
		return 1 + r
	}
	return height(root) >= 0
}
