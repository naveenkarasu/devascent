func digit_sum(n int) int {
	total := 0
	for n > 0 {
		total += n % 10
		n /= 10
	}
	return total
}
