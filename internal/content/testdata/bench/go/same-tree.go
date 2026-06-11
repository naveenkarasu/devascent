func is_same_tree(p *TreeNode, q *TreeNode) bool {
	if p == nil && q == nil {
		return true
	}
	if p == nil || q == nil || p.Val != q.Val {
		return false
	}
	return is_same_tree(p.Left, q.Left) && is_same_tree(p.Right, q.Right)
}
