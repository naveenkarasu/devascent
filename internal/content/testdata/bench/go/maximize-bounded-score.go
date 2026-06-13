func maximize_score(arr []int, k int) int {
	best := 0
	for _, x := range arr {
		diff := x - k
		if diff < 0 {
			diff = -diff
		}
		score := (2 * k) / (1 + 2*diff)
		if score > best {
			best = score
		}
	}
	return best
}
