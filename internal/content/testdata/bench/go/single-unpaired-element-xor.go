func find_unpaired(arr []int) int {
	result := 0
	for _, x := range arr {
		result ^= x
	}
	return result
}
