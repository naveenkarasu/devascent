func can_attend_meetings(intervals [][]int) bool {
	// sort intervals by start
	n := len(intervals)
	for i := 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			if intervals[j][0] < intervals[i][0] {
				intervals[i], intervals[j] = intervals[j], intervals[i]
			}
		}
	}
	for i := 1; i < n; i++ {
		if intervals[i][0] < intervals[i-1][1] {
			return false
		}
	}
	return true
}
