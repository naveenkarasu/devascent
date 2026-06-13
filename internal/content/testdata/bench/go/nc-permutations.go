import "sort"

func permute(nums []int) [][]int {
	res := [][]int{}
	var backtrack func(current, remaining []int)
	backtrack = func(current, remaining []int) {
		if len(remaining) == 0 {
			tmp := make([]int, len(current))
			copy(tmp, current)
			res = append(res, tmp)
			return
		}
		for i := range remaining {
			newRemaining := make([]int, 0, len(remaining)-1)
			newRemaining = append(newRemaining, remaining[:i]...)
			newRemaining = append(newRemaining, remaining[i+1:]...)
			backtrack(append(current, remaining[i]), newRemaining)
		}
	}
	backtrack([]int{}, nums)
	sort.Slice(res, func(i, j int) bool {
		a, b := res[i], res[j]
		for k := 0; k < len(a); k++ {
			if a[k] != b[k] {
				return a[k] < b[k]
			}
		}
		return false
	})
	return res
}
