func modular_exponentiation(x, n, m int) int {
	res := 1
	x = x % m
	for n > 0 {
		if n%2 != 0 {
			res = (res * x) % m
			n--
		} else {
			x = (x * x) % m
			n /= 2
		}
	}
	return res
}
