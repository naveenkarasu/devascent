func find_words(board [][]string, words []string) []string {
	// Build trie
	trie := map[string]interface{}{}
	for _, w := range words {
		node := trie
		for _, ch := range w {
			c := string(ch)
			if _, ok := node[c]; !ok {
				node[c] = map[string]interface{}{}
			}
			node = node[c].(map[string]interface{})
		}
		node["$"] = w
	}
	if len(board) == 0 || len(board[0]) == 0 {
		return []string{}
	}
	rows, cols := len(board), len(board[0])
	foundMap := map[string]bool{}
	var dfs func(r, c int, node map[string]interface{})
	dfs = func(r, c int, node map[string]interface{}) {
		ch := board[r][c]
		if ch == "#" {
			return
		}
		nxtIface, ok := node[ch]
		if !ok {
			return
		}
		nxt := nxtIface.(map[string]interface{})
		if w, ok := nxt["$"]; ok {
			foundMap[w.(string)] = true
		}
		board[r][c] = "#"
		dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
		for _, d := range dirs {
			nr, nc := r+d[0], c+d[1]
			if nr >= 0 && nr < rows && nc >= 0 && nc < cols && board[nr][nc] != "#" {
				dfs(nr, nc, nxt)
			}
		}
		board[r][c] = ch
	}
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			dfs(r, c, trie)
		}
	}
	result := []string{}
	for w := range foundMap {
		result = append(result, w)
	}
	// sort result
	for i := 0; i < len(result)-1; i++ {
		for j := i + 1; j < len(result); j++ {
			if result[j] < result[i] {
				result[i], result[j] = result[j], result[i]
			}
		}
	}
	return result
}
