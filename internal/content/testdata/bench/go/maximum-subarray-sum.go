func max_subarray_sum(nums []int) int {
	best := nums[0]
	current := nums[0]
	for _, x := range nums[1:] {
		if x > current+x {
			current = x
		} else {
			current = current + x
		}
		if current > best {
			best = current
		}
	}
	return best
}
