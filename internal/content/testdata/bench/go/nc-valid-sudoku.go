func is_valid_sudoku(board [][]string) bool {
	rows := make([]map[string]bool, 9)
	cols := make([]map[string]bool, 9)
	boxes := make([]map[string]bool, 9)
	for i := 0; i < 9; i++ {
		rows[i] = make(map[string]bool)
		cols[i] = make(map[string]bool)
		boxes[i] = make(map[string]bool)
	}
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			val := board[r][c]
			if val == "." {
				continue
			}
			boxIdx := (r/3)*3 + (c / 3)
			if rows[r][val] || cols[c][val] || boxes[boxIdx][val] {
				return false
			}
			rows[r][val] = true
			cols[c][val] = true
			boxes[boxIdx][val] = true
		}
	}
	return true
}
