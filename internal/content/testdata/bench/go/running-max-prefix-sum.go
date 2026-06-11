func peak_occupancy(initial int, deltas []int) int {
	current := initial
	peak := initial
	for _, d := range deltas {
		current += d
		if current > peak {
			peak = current
		}
	}
	return peak
}
