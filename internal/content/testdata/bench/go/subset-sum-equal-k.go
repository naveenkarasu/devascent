func subset_sum_to_k(nums []int, k int) bool {
	n := len(nums)
	prev := make([]bool, k+1)
	prev[0] = true
	if nums[0] <= k {
		prev[nums[0]] = true
	}
	for i := 1; i < n; i++ {
		curr := make([]bool, k+1)
		curr[0] = true
		for j := 1; j <= k; j++ {
			notTake := prev[j]
			take := false
			if nums[i] <= j {
				take = prev[j-nums[i]]
			}
			curr[j] = take || notTake
		}
		prev = curr
	}
	return prev[k]
}
