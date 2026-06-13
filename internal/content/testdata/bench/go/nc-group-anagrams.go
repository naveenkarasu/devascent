import (
	"sort"
)

func group_anagrams(strs []string) [][]string {
	groups := make(map[string][]string)
	for _, s := range strs {
		runes := []byte(s)
		sort.Slice(runes, func(i, j int) bool { return runes[i] < runes[j] })
		key := string(runes)
		groups[key] = append(groups[key], s)
	}
	res := make([][]string, 0, len(groups))
	for _, g := range groups {
		sorted := make([]string, len(g))
		copy(sorted, g)
		sort.Strings(sorted)
		res = append(res, sorted)
	}
	sort.Slice(res, func(i, j int) bool {
		for k := 0; k < len(res[i]) && k < len(res[j]); k++ {
			if res[i][k] != res[j][k] {
				return res[i][k] < res[j][k]
			}
		}
		return len(res[i]) < len(res[j])
	})
	return res
}
