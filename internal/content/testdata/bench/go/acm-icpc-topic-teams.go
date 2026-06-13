func best_team_coverage(topics [][]int) []int {
	n := len(topics)
	k := len(topics[0])
	maxTopics := 0
	bestCount := 0
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			covered := 0
			for t := 0; t < k; t++ {
				if topics[i][t] == 1 || topics[j][t] == 1 {
					covered++
				}
			}
			if covered > maxTopics {
				maxTopics = covered
				bestCount = 1
			} else if covered == maxTopics {
				bestCount++
			}
		}
	}
	return []int{maxTopics, bestCount}
}
