func match_colored_houses(left_colors []int, right_colors []int) []int {
	colorToPos := make(map[int]int)
	for i, c := range right_colors {
		colorToPos[c] = i + 1
	}
	result := make([]int, len(left_colors))
	for i, c := range left_colors {
		result[i] = colorToPos[c]
	}
	return result
}
