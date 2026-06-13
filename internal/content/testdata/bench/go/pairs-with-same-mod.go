func count_pairs_same_mod(nums []int, divisor int) int {
	freq := make(map[int]int)
	for _, x := range nums {
		freq[x%divisor]++
	}
	pairs := 0
	for _, cnt := range freq {
		pairs += cnt * (cnt - 1) / 2
	}
	return pairs
}
