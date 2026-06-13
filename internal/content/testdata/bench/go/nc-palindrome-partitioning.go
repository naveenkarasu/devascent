import "sort"

func partition(s string) [][]string {
	var res [][]string
	var backtrack func(start int, current []string)
	backtrack = func(start int, current []string) {
		if start == len(s) {
			tmp := make([]string, len(current))
			copy(tmp, current)
			res = append(res, tmp)
			return
		}
		for end := start + 1; end <= len(s); end++ {
			substr := s[start:end]
			if isPalin(substr) {
				current = append(current, substr)
				backtrack(end, current)
				current = current[:len(current)-1]
			}
		}
	}
	backtrack(0, []string{})
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

func isPalin(t string) bool {
	for i, j := 0, len(t)-1; i < j; i, j = i+1, j-1 {
		if t[i] != t[j] {
			return false
		}
	}
	return true
}
