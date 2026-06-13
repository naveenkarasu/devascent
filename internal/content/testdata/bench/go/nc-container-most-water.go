func max_area(heights []int) int {
	i, j := 0, len(heights)-1
	best := 0
	for i < j {
		h := heights[i]
		if heights[j] < h {
			h = heights[j]
		}
		area := h * (j - i)
		if area > best {
			best = area
		}
		if heights[i] < heights[j] {
			i++
		} else {
			j--
		}
	}
	return best
}
