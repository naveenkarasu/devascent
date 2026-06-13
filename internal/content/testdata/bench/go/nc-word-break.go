func word_break(s string, word_dict []string) bool {
	wordSet := make(map[string]bool)
	for _, w := range word_dict {
		wordSet[w] = true
	}
	n := len(s)
	dp := make([]bool, n+1)
	dp[0] = true
	for i := 1; i <= n; i++ {
		for j := 0; j < i; j++ {
			if dp[j] && wordSet[s[j:i]] {
				dp[i] = true
				break
			}
		}
	}
	return dp[n]
}
