func is_armstrong_3digit(n int) bool {
	if n < 100 || n > 999 {
		return false
	}
	s := 0
	temp := n
	for temp > 0 {
		r := temp % 10
		s += r * r * r
		temp /= 10
	}
	return s == n
}
