import "sort"

func merge_intervals(intervals [][]int) [][]int {
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})
	res := [][]int{}
	for _, iv := range intervals {
		s, e := iv[0], iv[1]
		if len(res) > 0 && s <= res[len(res)-1][1] {
			if e > res[len(res)-1][1] {
				res[len(res)-1][1] = e
			}
		} else {
			res = append(res, []int{s, e})
		}
	}
	return res
}
