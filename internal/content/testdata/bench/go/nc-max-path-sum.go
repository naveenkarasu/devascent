func max_path_sum(root *TreeNode) int {
	best := -(1 << 62)
	var gain func(node *TreeNode) int
	gain = func(node *TreeNode) int {
		if node == nil {
			return 0
		}
		l := gain(node.Left)
		if l < 0 {
			l = 0
		}
		r := gain(node.Right)
		if r < 0 {
			r = 0
		}
		sum := node.Val + l + r
		if sum > best {
			best = sum
		}
		if l > r {
			return node.Val + l
		}
		return node.Val + r
	}
	gain(root)
	return best
}
