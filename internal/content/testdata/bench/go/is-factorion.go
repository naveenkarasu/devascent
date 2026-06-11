func is_factorion(n int) bool {
	factorial := func(x int) int {
		result := 1
		for i := 2; i <= x; i++ {
			result *= i
		}
		return result
	}
	original := n
	s := 0
	for n > 0 {
		d := n % 10
		s += factorial(d)
		n /= 10
	}
	return s == original
}
