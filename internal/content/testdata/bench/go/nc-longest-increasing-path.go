func longest_increasing_path(matrix [][]int) int {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return 0
	}
	m, n := len(matrix), len(matrix[0])
	memo := make([][]int, m)
	for i := range memo {
		memo[i] = make([]int, n)
	}

	dirs := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	var dfs func(r, c int) int
	dfs = func(r, c int) int {
		if memo[r][c] != 0 {
			return memo[r][c]
		}
		best := 1
		for _, d := range dirs {
			nr, nc := r+d[0], c+d[1]
			if nr >= 0 && nr < m && nc >= 0 && nc < n && matrix[nr][nc] > matrix[r][c] {
				v := 1 + dfs(nr, nc)
				if v > best {
					best = v
				}
			}
		}
		memo[r][c] = best
		return best
	}

	res := 0
	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			v := dfs(r, c)
			if v > res {
				res = v
			}
		}
	}
	return res
}
