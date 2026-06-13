func can_complete_circuit(gas []int, cost []int) int {
	totalGas, totalCost := 0, 0
	for i := range gas {
		totalGas += gas[i]
		totalCost += cost[i]
	}
	if totalGas < totalCost {
		return -1
	}
	total := 0
	start := 0
	for i := range gas {
		total += gas[i] - cost[i]
		if total < 0 {
			start = i + 1
			total = 0
		}
	}
	return start
}
