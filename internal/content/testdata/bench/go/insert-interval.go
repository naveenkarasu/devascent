func insert_interval(intervals [][]int, newInterval []int) [][]int {
	res := [][]int{}
	i, n := 0, len(intervals)
	s, e := newInterval[0], newInterval[1]
	for i < n && intervals[i][1] < s {
		res = append(res, intervals[i])
		i++
	}
	for i < n && intervals[i][0] <= e {
		if intervals[i][0] < s {
			s = intervals[i][0]
		}
		if intervals[i][1] > e {
			e = intervals[i][1]
		}
		i++
	}
	res = append(res, []int{s, e})
	for i < n {
		res = append(res, intervals[i])
		i++
	}
	return res
}
