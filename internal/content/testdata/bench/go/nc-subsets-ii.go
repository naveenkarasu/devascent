import "sort"

func subsets_with_dup(nums []int) [][]int {
	sort.Ints(nums)
	var res [][]int
	var backtrack func(start int, current []int)
	backtrack = func(start int, current []int) {
		tmp := make([]int, len(current))
		copy(tmp, current)
		res = append(res, tmp)
		for i := start; i < len(nums); i++ {
			if i > start && nums[i] == nums[i-1] {
				continue
			}
			current = append(current, nums[i])
			backtrack(i+1, current)
			current = current[:len(current)-1]
		}
	}
	backtrack(0, []int{})
	sort.Slice(res, func(i, j int) bool {
		a, b := res[i], res[j]
		for k := 0; k < len(a) && k < len(b); k++ {
			if a[k] != b[k] {
				return a[k] < b[k]
			}
		}
		return len(a) < len(b)
	})
	return res
}
