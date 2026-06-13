func min_cost_climbing_stairs(cost []int) int {
	n := len(cost)
	for i := 2; i < n; i++ {
		cost[i] += min(cost[i-1], cost[i-2])
	}
	return min(cost[n-1], cost[n-2])
}
