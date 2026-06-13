import "sort"

func combination_sum(candidates []int, target int) [][]int {
	sort.Ints(candidates)
	res := [][]int{}
	var backtrack func(start, remaining int, current []int)
	backtrack = func(start, remaining int, current []int) {
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
			current = append(current, candidates[i])
			backtrack(i, remaining-candidates[i], current)
			current = current[:len(current)-1]
		}
	}
	backtrack(0, target, []int{})
	for _, c := range res {
		sort.Ints(c)
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
