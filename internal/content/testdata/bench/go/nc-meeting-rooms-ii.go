func min_meeting_rooms(intervals [][]int) int {
	if len(intervals) == 0 {
		return 0
	}
	n := len(intervals)
	starts := make([]int, n)
	ends := make([]int, n)
	for i, iv := range intervals {
		starts[i] = iv[0]
		ends[i] = iv[1]
	}
	// sort both (simple insertion sort)
	for i := 1; i < n; i++ {
		for j := i; j > 0 && starts[j] < starts[j-1]; j-- {
			starts[j], starts[j-1] = starts[j-1], starts[j]
		}
	}
	for i := 1; i < n; i++ {
		for j := i; j > 0 && ends[j] < ends[j-1]; j-- {
			ends[j], ends[j-1] = ends[j-1], ends[j]
		}
	}
	rooms, best := 0, 0
	s, e := 0, 0
	for s < n {
		if starts[s] < ends[e] {
			rooms++
			s++
			if rooms > best {
				best = rooms
			}
		} else {
			rooms--
			e++
		}
	}
	return best
}
