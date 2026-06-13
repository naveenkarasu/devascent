func last_remaining(n int) int {
	if n == 1 {
		return 1
	}
	return 2 * (1 + n/2 - last_remaining(n/2))
}
