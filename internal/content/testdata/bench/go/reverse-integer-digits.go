func reverse_integer(n int) int {
	reversed := 0
	for n > 0 {
		reversed = reversed*10 + n%10
		n /= 10
	}
	return reversed
}
