func max_sliding_window(nums []int, k int) []int {
	n := len(nums)
	result := []int{}
	dq := []int{} // indices, decreasing by value
	for i := 0; i < n; i++ {
		for len(dq) > 0 && dq[0] < i-k+1 {
			dq = dq[1:]
		}
		for len(dq) > 0 && nums[dq[len(dq)-1]] < nums[i] {
			dq = dq[:len(dq)-1]
		}
		dq = append(dq, i)
		if i >= k-1 {
			result = append(result, nums[dq[0]])
		}
	}
	return result
}
