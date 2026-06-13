func find_target_sum_ways(nums []int, target int) int {
	dp := map[int]int{0: 1}
	for _, num := range nums {
		nextDp := make(map[int]int)
		for s, cnt := range dp {
			nextDp[s+num] += cnt
			nextDp[s-num] += cnt
		}
		dp = nextDp
	}
	return dp[target]
}
