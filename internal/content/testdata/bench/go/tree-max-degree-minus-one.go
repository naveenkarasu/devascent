func trap_node(n int, edges [][]int) int {
	degree := make([]int, n)
	for _, e := range edges {
		degree[e[0]]++
		degree[e[1]]++
	}
	maxDeg := 0
	for _, d := range degree {
		if d > maxDeg {
			maxDeg = d
		}
	}
	return maxDeg - 1
}
