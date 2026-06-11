func can_partition(nums []int) bool {
	total := 0
	for _, n := range nums {
		total += n
	}
	if total%2 != 0 {
		return false
	}
	target := total / 2
	dp := make([]bool, target+1)
	dp[0] = true
	for _, n := range nums {
		for s := target; s >= n; s-- {
			if dp[s-n] {
				dp[s] = true
			}
		}
		if dp[target] {
			return true
		}
	}
	return false
}
