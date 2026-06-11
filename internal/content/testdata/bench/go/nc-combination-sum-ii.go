import "sort"

func combination_sum2(candidates []int, target int) [][]int {
	sort.Ints(candidates)
	var res [][]int
	var backtrack func(start int, current []int, remaining int)
	backtrack = func(start int, current []int, remaining int) {
		if remaining == 0 {
			tmp := make([]int, len(current))
			copy(tmp, current)
			res = append(res, tmp)
			return
		}
		for i := start; i < len(candidates); i++ {
			if candidates[i] > remaining {
				break
			}
			if i > start && candidates[i] == candidates[i-1] {
				continue
			}
			current = append(current, candidates[i])
			backtrack(i+1, current, remaining-candidates[i])
			current = current[:len(current)-1]
		}
	}
	backtrack(0, []int{}, target)
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
