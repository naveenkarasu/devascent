func trap(heights []int) int {
	if len(heights) == 0 {
		return 0
	}
	i, j := 0, len(heights)-1
	leftMax, rightMax := heights[i], heights[j]
	total := 0
	for i < j {
		if leftMax < rightMax {
			i++
			if heights[i] > leftMax {
				leftMax = heights[i]
			}
			total += leftMax - heights[i]
		} else {
			j--
			if heights[j] > rightMax {
				rightMax = heights[j]
			}
			total += rightMax - heights[j]
		}
	}
	return total
}
