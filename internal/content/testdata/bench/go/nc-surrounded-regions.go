func solve(board [][]string) [][]string {
	if len(board) == 0 || len(board[0]) == 0 {
		return board
	}
	m, n := len(board), len(board[0])
	safe := make([][]bool, m)
	for i := range safe {
		safe[i] = make([]bool, n)
	}
	var dfs func(i, j int)
	dfs = func(i, j int) {
		if i < 0 || j < 0 || i >= m || j >= n {
			return
		}
		if board[i][j] != "O" || safe[i][j] {
			return
		}
		safe[i][j] = true
		dfs(i+1, j)
		dfs(i-1, j)
		dfs(i, j+1)
		dfs(i, j-1)
	}
	for i := 0; i < m; i++ {
		dfs(i, 0)
		dfs(i, n-1)
	}
	for j := 0; j < n; j++ {
		dfs(0, j)
		dfs(m-1, j)
	}
	result := make([][]string, m)
	for i := 0; i < m; i++ {
		result[i] = make([]string, n)
		for j := 0; j < n; j++ {
			if board[i][j] == "O" && !safe[i][j] {
				result[i][j] = "X"
			} else {
				result[i][j] = board[i][j]
			}
		}
	}
	return result
}
