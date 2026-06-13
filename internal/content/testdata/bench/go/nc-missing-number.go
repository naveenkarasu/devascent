func missing_number(nums []int) int {
	n := len(nums)
	expected := n * (n + 1) / 2
	sum := 0
	for _, v := range nums {
		sum += v
	}
	return expected - sum
}
