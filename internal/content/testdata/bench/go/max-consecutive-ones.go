func max_consecutive_ones(nums []int) int {
	maxRun := 0
	current := 0
	for _, x := range nums {
		if x == 1 {
			current++
			if current > maxRun {
				maxRun = current
			}
		} else {
			current = 0
		}
	}
	return maxRun
}
