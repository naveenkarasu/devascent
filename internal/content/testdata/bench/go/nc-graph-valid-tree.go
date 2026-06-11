func valid_tree(n int, edges [][]int) bool {
	if len(edges) != n-1 {
		return false
	}
	adj := make([][]int, n)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	visited := make([]bool, n)
	var dfs func(node int)
	dfs = func(node int) {
		visited[node] = true
		for _, nb := range adj[node] {
			if !visited[nb] {
				dfs(nb)
			}
		}
	}
	dfs(0)
	for _, v := range visited {
		if !v {
			return false
		}
	}
	return true
}
