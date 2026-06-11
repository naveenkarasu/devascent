func modular_sqrt(n int, m int) int {
	for x := 0; x < m; x++ {
		if (x*x)%m == n {
			return x
		}
	}
	return -1
}
