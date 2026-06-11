func find_unique_among_k(k int, arr []int) int {
	counts := make(map[int]int)
	for _, v := range arr {
		counts[v]++
	}
	for _, v := range arr {
		if counts[v]%k != 0 {
			return v
		}
	}
	return -1
}
