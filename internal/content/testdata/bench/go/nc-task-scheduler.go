func least_interval(tasks []string, n int) int {
	counts := map[string]int{}
	for _, t := range tasks {
		counts[t]++
	}
	maxCount := 0
	for _, c := range counts {
		if c > maxCount {
			maxCount = c
		}
	}
	numMax := 0
	for _, c := range counts {
		if c == maxCount {
			numMax++
		}
	}
	result := (maxCount-1)*(n+1) + numMax
	if len(tasks) > result {
		return len(tasks)
	}
	return result
}
