func clear_full_rows(board [][]string, empty string) [][]string {
	if len(board) == 0 {
		return board
	}
	cols := len(board[0])
	var surviving [][]string
	for _, row := range board {
		hasEmpty := false
		for _, cell := range row {
			if cell == empty {
				hasEmpty = true
				break
			}
		}
		if hasEmpty {
			surviving = append(surviving, row)
		}
	}
	clearedCount := len(board) - len(surviving)
	result := make([][]string, clearedCount)
	for i := range result {
		blankRow := make([]string, cols)
		for j := range blankRow {
			blankRow[j] = empty
		}
		result[i] = blankRow
	}
	result = append(result, surviving...)
	return result
}
