func max_coins(nums []int) int {
	nums = append([]int{1}, append(nums, 1)...)
	n := len(nums)
	dp := make([][]int, n)
	for i := range dp {
		dp[i] = make([]int, n)
	}
	for length := 2; length < n; length++ {
		for left := 0; left < n-length; left++ {
			right := left + length
			for k := left + 1; k < right; k++ {
				coins := nums[left]*nums[k]*nums[right] + dp[left][k] + dp[k][right]
				if coins > dp[left][right] {
					dp[left][right] = coins
				}
			}
		}
	}
	return dp[0][n-1]
}
