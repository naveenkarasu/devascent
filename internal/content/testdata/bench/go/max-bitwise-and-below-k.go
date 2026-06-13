func max_and_below_k(n int, k int) int {
	if (k-1)|k <= n {
		return k - 1
	}
	return k - 2
}
