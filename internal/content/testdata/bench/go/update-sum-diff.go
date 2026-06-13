func update_pair(a int, b int) []int {
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return []int{a + b, diff}
}
