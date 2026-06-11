func triple_bucket_sum(counts []int) int {
	const MOD = 1_000_000_007
	ans := 0
	n := len(counts)
	for i := 0; i < n-2; i++ {
		for j := i + 1; j < n-1; j++ {
			for k := j + 1; k < n; k++ {
				ans = (ans + counts[i]*counts[j]*counts[k]) % MOD
			}
		}
	}
	return ans
}
