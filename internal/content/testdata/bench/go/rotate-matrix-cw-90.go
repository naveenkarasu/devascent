func rotate_matrix_cw(matrix [][]int) [][]int {
	n := len(matrix)
	result := make([][]int, n)
	for i := range result {
		result[i] = make([]int, n)
		for j := 0; j < n; j++ {
			result[i][j] = matrix[n-1-j][i]
		}
	}
	return result
}
