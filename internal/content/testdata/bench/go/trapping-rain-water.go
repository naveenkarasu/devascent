func trap_water(height []int) int {
	l, r := 0, len(height)-1
	level, water := 0, 0
	for l < r {
		lower := height[l]
		if height[r] < lower {
			lower = height[r]
		}
		if height[l] < height[r] {
			l++
		} else {
			r--
		}
		if lower > level {
			level = lower
		}
		water += level - lower
	}
	return water
}
