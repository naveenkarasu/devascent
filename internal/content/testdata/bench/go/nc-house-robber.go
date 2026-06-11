func rob(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	if len(nums) == 1 {
		return nums[0]
	}
	prev2, prev1 := nums[0], max(nums[0], nums[1])
	for i := 2; i < len(nums); i++ {
		curr := max(prev1, prev2+nums[i])
		prev2, prev1 = prev1, curr
	}
	return prev1
}
