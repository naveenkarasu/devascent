import "sort"

type graphNode struct {
	Val       int
	Neighbors []*graphNode
}

func cloneNode(node *graphNode, mp map[int]*graphNode) *graphNode {
	if node == nil {
		return nil
	}
	if copy, ok := mp[node.Val]; ok {
		return copy
	}
	cp := &graphNode{Val: node.Val}
	mp[node.Val] = cp
	for _, nb := range node.Neighbors {
		cp.Neighbors = append(cp.Neighbors, cloneNode(nb, mp))
	}
	return cp
}

func clone_graph(adj [][]int) [][]int {
	if len(adj) == 0 {
		return [][]int{}
	}
	nodes := make(map[int]*graphNode)
	for i := range adj {
		nodes[i+1] = &graphNode{Val: i + 1}
	}
	for i, nbrs := range adj {
		for _, x := range nbrs {
			nodes[i+1].Neighbors = append(nodes[i+1].Neighbors, nodes[x])
		}
	}
	mp := make(map[int]*graphNode)
	cloned := cloneNode(nodes[1], mp)

	out := make([][]int, len(adj))
	seen := make(map[int]bool)
	stack := []*graphNode{cloned}
	for len(stack) > 0 {
		nd := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if seen[nd.Val] {
			continue
		}
		seen[nd.Val] = true
		vals := []int{}
		for _, n := range nd.Neighbors {
			vals = append(vals, n.Val)
			stack = append(stack, n)
		}
		sort.Ints(vals)
		out[nd.Val-1] = vals
	}
	return out
}
