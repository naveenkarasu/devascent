func sequence_equation(p []int) []int {
	n := len(p)
	inv := make([]int, n+1)
	for i := 0; i < n; i++ {
		inv[p[i]] = i + 1
	}
	result := make([]int, n)
	for x := 1; x <= n; x++ {
		result[x-1] = inv[inv[x]]
	}
	return result
}
