func find_celebrity(knows_matrix [][]bool) int {
	n := len(knows_matrix)
	candidate := 0
	for i := 1; i < n; i++ {
		if knows_matrix[candidate][i] {
			candidate = i
		}
	}
	for i := 0; i < n; i++ {
		if i != candidate {
			if knows_matrix[candidate][i] || !knows_matrix[i][candidate] {
				return -1
			}
		}
	}
	return candidate
}
