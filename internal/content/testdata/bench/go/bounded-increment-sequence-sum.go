func bounded_sequence_sum(n int, b int, x int, y int) int {
	a := make([]int, n+1)
	a[0] = 0
	for i := 1; i <= n; i++ {
		if a[i-1]+x <= b {
			a[i] = a[i-1] + x
		} else {
			a[i] = a[i-1] - y
		}
	}
	sum := 0
	for _, v := range a {
		sum += v
	}
	return sum
}
