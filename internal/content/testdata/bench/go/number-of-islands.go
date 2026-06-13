func num_islands(grid [][]int) int {
	if len(grid) == 0 {
		return 0
	}
	m, n := len(grid), len(grid[0])
	seen := make([][]bool, m)
	for i := range seen {
		seen[i] = make([]bool, n)
	}
	var dfs func(i, j int)
	dfs = func(i, j int) {
		if i < 0 || j < 0 || i >= m || j >= n {
			return
		}
		if grid[i][j] == 0 || seen[i][j] {
			return
		}
		seen[i][j] = true
		dfs(i+1, j)
		dfs(i-1, j)
		dfs(i, j+1)
		dfs(i, j-1)
	}
	count := 0
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if grid[i][j] == 1 && !seen[i][j] {
				count++
				dfs(i, j)
			}
		}
	}
	return count
}
