func exist(board [][]string, word string) bool {
	rows := len(board)
	cols := len(board[0])
	visited := make([][]bool, rows)
	for i := range visited {
		visited[i] = make([]bool, cols)
	}
	runes := []rune(word)
	var dfs func(r, c, index int) bool
	dfs = func(r, c, index int) bool {
		if index == len(runes) {
			return true
		}
		if r < 0 || r >= rows || c < 0 || c >= cols {
			return false
		}
		if visited[r][c] {
			return false
		}
		if board[r][c] != string(runes[index]) {
			return false
		}
		visited[r][c] = true
		found := dfs(r+1, c, index+1) || dfs(r-1, c, index+1) ||
			dfs(r, c+1, index+1) || dfs(r, c-1, index+1)
		visited[r][c] = false
		return found
	}
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if dfs(r, c, 0) {
				return true
			}
		}
	}
	return false
}
