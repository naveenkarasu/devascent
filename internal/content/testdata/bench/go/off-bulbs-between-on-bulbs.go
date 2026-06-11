func off_bulbs_between_lit(bulbs []int) int {
	n := len(bulbs)
	firstOn := -1
	lastOn := -1
	for i := 0; i < n; i++ {
		if bulbs[i] == 1 {
			firstOn = i
			break
		}
	}
	for i := n - 1; i >= 0; i-- {
		if bulbs[i] == 1 {
			lastOn = i
			break
		}
	}
	if firstOn == -1 {
		return 0
	}
	count := 0
	for j := firstOn; j <= lastOn; j++ {
		if bulbs[j] == 0 {
			count++
		}
	}
	return count
}
