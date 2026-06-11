func find_redundant_connection(edges [][]int) []int {
	n := len(edges)
	parent := make([]int, n+1)
	rank := make([]int, n+1)
	for i := range parent {
		parent[i] = i
	}

	var find func(x int) int
	find = func(x int) int {
		for parent[x] != x {
			parent[x] = parent[parent[x]]
			x = parent[x]
		}
		return x
	}

	union := func(a, b int) bool {
		pa, pb := find(a), find(b)
		if pa == pb {
			return false
		}
		if rank[pa] < rank[pb] {
			pa, pb = pb, pa
		}
		parent[pb] = pa
		if rank[pa] == rank[pb] {
			rank[pa]++
		}
		return true
	}

	for _, e := range edges {
		if !union(e[0], e[1]) {
			return []int{e[0], e[1]}
		}
	}
	return []int{}
}
