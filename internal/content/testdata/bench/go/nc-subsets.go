import "sort"

func subsets(nums []int) [][]int {
	res := [][]int{}
	var backtrack func(start int, current []int)
	backtrack = func(start int, current []int) {
		tmp := make([]int, len(current))
		copy(tmp, current)
		res = append(res, tmp)
		for i := start; i < len(nums); i++ {
			current = append(current, nums[i])
			backtrack(i+1, current)
			current = current[:len(current)-1]
		}
	}
	backtrack(0, []int{})
	for _, s := range res {
		sort.Ints(s)
	}
	sort.Slice(res, func(i, j int) bool {
		a, b := res[i], res[j]
		minLen := len(a)
		if len(b) < minLen {
			minLen = len(b)
		}
		for k := 0; k < minLen; k++ {
			if a[k] != b[k] {
				return a[k] < b[k]
			}
		}
		return len(a) < len(b)
	})
	return res
}
