func rectangle_perimeter_area(a, b int) []int {
	if a <= 0 || b <= 0 {
		return []int{0, 0}
	}
	return []int{(a + b) * 2, a * b}
}
