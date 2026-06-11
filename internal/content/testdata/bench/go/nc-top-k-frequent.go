import "sort"

func top_k_frequent(nums []int, k int) []int {
	freq := make(map[int]int)
	for _, n := range nums {
		freq[n]++
	}
	keys := make([]int, 0, len(freq))
	for key := range freq {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		if freq[keys[i]] != freq[keys[j]] {
			return freq[keys[i]] > freq[keys[j]]
		}
		return keys[i] < keys[j]
	})
	top := keys[:k]
	sort.Ints(top)
	return top
}
