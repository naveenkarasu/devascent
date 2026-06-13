func max_area_of_island(grid [][]int) int {
	if len(grid) == 0 {
		return 0
	}
	m, n := len(grid), len(grid[0])
	seen := make([][]bool, m)
	for i := range seen {
		seen[i] = make([]bool, n)
	}
	var dfs func(i, j int) int
	dfs = func(i, j int) int {
		if i < 0 || j < 0 || i >= m || j >= n {
			return 0
		}
		if grid[i][j] == 0 || seen[i][j] {
			return 0
		}
		seen[i][j] = true
		return 1 + dfs(i+1, j) + dfs(i-1, j) + dfs(i, j+1) + dfs(i, j-1)
	}
	best := 0
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if grid[i][j] == 1 && !seen[i][j] {
				area := dfs(i, j)
				if area > best {
					best = area
				}
			}
		}
	}
	return best
}
