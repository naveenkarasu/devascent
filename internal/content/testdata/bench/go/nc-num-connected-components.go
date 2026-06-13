func count_components(n int, edges [][]int) int {
	parent := make([]int, n)
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
	union := func(a, b int) int {
		pa, pb := find(a), find(b)
		if pa == pb {
			return 0
		}
		parent[pa] = pb
		return 1
	}
	merged := 0
	for _, e := range edges {
		merged += union(e[0], e[1])
	}
	return n - merged
}
