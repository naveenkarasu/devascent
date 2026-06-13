func count_modeq_pairs(n int, m int) int {
	count := 0
	for a := 1; a <= n; a++ {
		for b := a + 1; b <= n; b++ {
			if (m%a)%b == (m%b)%a {
				count++
			}
		}
	}
	return count
}
