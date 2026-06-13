func min_eating_speed(piles []int, h int) int {
	lo, hi := 1, 0
	for _, p := range piles {
		if p > hi {
			hi = p
		}
	}
	ans := hi
	for lo <= hi {
		mid := (lo + hi) / 2
		hours := 0
		for _, p := range piles {
			hours += (p + mid - 1) / mid
		}
		if hours <= h {
			ans = mid
			hi = mid - 1
		} else {
			lo = mid + 1
		}
	}
	return ans
}
