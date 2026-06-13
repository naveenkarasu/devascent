import "sort"

func set_union(a, b []int) []int {
	seen := map[int]bool{}
	for _, v := range a {
		seen[v] = true
	}
	for _, v := range b {
		seen[v] = true
	}
	result := make([]int, 0, len(seen))
	for k := range seen {
		result = append(result, k)
	}
	sort.Ints(result)
	return result
}
