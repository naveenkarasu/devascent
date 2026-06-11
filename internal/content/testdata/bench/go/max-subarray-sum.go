func max_subarray(nums []int) int {
	current := nums[0]
	best := nums[0]
	for i := 1; i < len(nums); i++ {
		if nums[i] > current+nums[i] {
			current = nums[i]
		} else {
			current = current + nums[i]
		}
		if current > best {
			best = current
		}
	}
	return best
}
