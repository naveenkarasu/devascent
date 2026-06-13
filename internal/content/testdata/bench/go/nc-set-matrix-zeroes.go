func set_zeroes(matrix [][]int) [][]int {
	rows := map[int]bool{}
	cols := map[int]bool{}
	for r := range matrix {
		for c := range matrix[r] {
			if matrix[r][c] == 0 {
				rows[r] = true
				cols[c] = true
			}
		}
	}
	for r := range matrix {
		for c := range matrix[r] {
			if rows[r] || cols[c] {
				matrix[r][c] = 0
			}
		}
	}
	return matrix
}
