func all_permutations(nums []int) [][]int {
	if len(nums) == 0 {
		return [][]int{{}}
	}
	result := [][]int{}
	for i, n := range nums {
		rest := make([]int, 0, len(nums)-1)
		rest = append(rest, nums[:i]...)
		rest = append(rest, nums[i+1:]...)
		for _, perm := range all_permutations(rest) {
			combo := make([]int, 0, len(perm)+1)
			combo = append(combo, n)
			combo = append(combo, perm...)
			result = append(result, combo)
		}
	}
	return result
}
