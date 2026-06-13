func free_vertices(n int, m int) int {
	k := 0
	for k*(k-1)/2 < m {
		k++
	}
	return n - k
}
