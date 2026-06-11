func two_sum(nums []int, target int) []int {
	seen := map[int]int{}
	for i, n := range nums {
		need := target - n
		if j, ok := seen[need]; ok {
			return []int{j, i}
		}
		seen[n] = i
	}
	return []int{}
}
