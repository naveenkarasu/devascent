func sum_squared_digits(n int) int {
	total := 0
	for n > 0 {
		d := n % 10
		total += d * d
		n /= 10
	}
	return total
}
