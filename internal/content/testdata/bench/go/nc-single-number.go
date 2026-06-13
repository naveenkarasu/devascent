func single_number(nums []int) int {
	result := 0
	for _, n := range nums {
		result ^= n
	}
	return result
}
