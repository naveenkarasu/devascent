func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func smallest_trailing_zero_multiple(a int, b int) int {
	power := 1
	for i := 0; i < b; i++ {
		power *= 10
	}
	return (a / gcd(a, power)) * power
}
