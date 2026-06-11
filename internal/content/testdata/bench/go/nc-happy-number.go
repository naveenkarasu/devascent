func is_happy(n int) bool {
	seen := map[int]bool{}
	for n != 1 {
		if seen[n] {
			return false
		}
		seen[n] = true
		total := 0
		for n > 0 {
			d := n % 10
			total += d * d
			n /= 10
		}
		n = total
	}
	return true
}
