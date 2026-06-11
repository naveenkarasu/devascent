func build_tree(preorder []int, inorder []int) *TreeNode {
	idx := map[int]int{}
	for i, v := range inorder {
		idx[v] = i
	}
	pre := 0
	var build func(lo, hi int) *TreeNode
	build = func(lo, hi int) *TreeNode {
		if lo > hi {
			return nil
		}
		val := preorder[pre]
		pre++
		node := &TreeNode{Val: val}
		mid := idx[val]
		node.Left = build(lo, mid-1)
		node.Right = build(mid+1, hi)
		return node
	}
	return build(0, len(inorder)-1)
}
