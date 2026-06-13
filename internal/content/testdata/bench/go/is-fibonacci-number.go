func is_fibonacci(n int) bool {
	if n <= 0 {
		return false
	}
	a, b := 1, 1
	for a < n {
		a, b = b, a+b
	}
	return a == n
}
