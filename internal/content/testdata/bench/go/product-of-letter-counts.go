func product_of_counts(s string) int {
	const MOD = 1_000_000_007
	freq := make(map[rune]int)
	for _, ch := range s {
		freq[ch]++
	}
	result := 1
	for _, cnt := range freq {
		result = (result * cnt) % MOD
	}
	return result
}
