func my_pow(x int, n int) int {
	if n == 0 {
		return 1
	}
	if n < 0 {
		// For integer base with negative exponent, result would be fractional
		// but per prompt tests use non-negative exponents
		return 0
	}
	half := my_pow(x, n/2)
	if n%2 == 0 {
		return half * half
	}
	return half * half * x
}
