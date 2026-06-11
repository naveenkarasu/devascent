func rob_circular(nums []int) int {
	if len(nums) == 1 {
		return nums[0]
	}
	return max(robLinear(nums[:len(nums)-1]), robLinear(nums[1:]))
}

func robLinear(houses []int) int {
	if len(houses) == 0 {
		return 0
	}
	if len(houses) == 1 {
		return houses[0]
	}
	prev2, prev1 := houses[0], max(houses[0], houses[1])
	for i := 2; i < len(houses); i++ {
		curr := max(prev1, prev2+houses[i])
		prev2, prev1 = prev1, curr
	}
	return prev1
}
