func max_div_plus_mod(l int, r int, a int) int {
	best := r/a + r%a
	multiple := (r / a) * a
	candidate := multiple - 1
	if l <= candidate && candidate <= r {
		val := candidate/a + candidate%a
		if val > best {
			best = val
		}
	}
	return best
}
