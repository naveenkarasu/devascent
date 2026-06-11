func top_two_sum(nums []int) int {
	first, second := 0, 0
	for _, x := range nums {
		if x >= first {
			second = first
			first = x
		} else if x > second {
			second = x
		}
	}
	return first + second
}
