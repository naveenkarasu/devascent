func min_distance(word1 string, word2 string) int {
	m, n := len(word1), len(word2)
	dp := make([]int, n+1)
	for j := 0; j <= n; j++ {
		dp[j] = j
	}
	for i := 1; i <= m; i++ {
		prev := dp[0]
		dp[0] = i
		for j := 1; j <= n; j++ {
			temp := dp[j]
			if word1[i-1] == word2[j-1] {
				dp[j] = prev
			} else {
				minVal := prev
				if dp[j] < minVal {
					minVal = dp[j]
				}
				if dp[j-1] < minVal {
					minVal = dp[j-1]
				}
				dp[j] = 1 + minVal
			}
			prev = temp
		}
	}
	return dp[n]
}
