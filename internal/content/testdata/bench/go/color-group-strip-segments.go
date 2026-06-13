func count_segments(groups [][]int, strip []int) int {
	colorToGroup := make(map[int]int)
	for idx, group := range groups {
		for _, color := range group {
			colorToGroup[color] = idx
		}
	}
	if len(strip) == 0 {
		return 0
	}
	segments := 1
	for i := 1; i < len(strip); i++ {
		g1, ok1 := colorToGroup[strip[i-1]]
		g2, ok2 := colorToGroup[strip[i]]
		if !ok1 || !ok2 || g1 != g2 {
			segments++
		}
	}
	return segments
}
