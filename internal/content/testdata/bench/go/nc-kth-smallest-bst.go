func kth_smallest(root *TreeNode, k int) int {
	stack := []*TreeNode{}
	cur := root
	for len(stack) > 0 || cur != nil {
		for cur != nil {
			stack = append(stack, cur)
			cur = cur.Left
		}
		cur = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		k--
		if k == 0 {
			return cur.Val
		}
		cur = cur.Right
	}
	return -1
}
